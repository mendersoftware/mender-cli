// Copyright 2021 Northern.tech AS
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package cmd

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"sync"

	"github.com/google/uuid"
	"github.com/mendersoftware/go-lib-micro/ws"
	"github.com/mendersoftware/go-lib-micro/ws/portforward"
	wspf "github.com/mendersoftware/go-lib-micro/ws/portforward"
	"github.com/vmihailenco/msgpack"
)

const portForwardUDPChannelSize = 20

type UDPPortForwarder struct {
	conn          *net.UDPConn
	remoteHost    string
	remotePort    uint16
	sourceAddr    *net.UDPAddr
	waitGroupAcks *sync.WaitGroup
}

func NewUDPPortForwarder(bindingHost string, localPort uint16, remoteHost string, remotePort uint16) (*UDPPortForwarder, error) {
	fmt.Printf("Forwarding from udp/%s:%d -> udp/%s:%d\n", bindingHost, localPort, remoteHost, remotePort)
	sAddr, err := net.ResolveUDPAddr("udp", bindingHost+":"+strconv.Itoa(int(localPort)))
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP(protocolUDP, sAddr)
	if err != nil {
		return nil, err
	}
	return &UDPPortForwarder{
		conn:          conn,
		remoteHost:    remoteHost,
		remotePort:    remotePort,
		waitGroupAcks: &sync.WaitGroup{},
	}, nil
}

func (p *UDPPortForwarder) Run(ctx context.Context, sessionID string, msgChan chan *ws.ProtoMsg, recvChans map[string]chan *ws.ProtoMsg) {
	// listen for new connections
	defer p.conn.Close()

	connectionUUID, _ := uuid.NewUUID()
	connectionID := connectionUUID.String()
	recvChan := make(chan *ws.ProtoMsg, portForwardUDPChannelSize)
	recvChans[connectionID] = recvChan

	protocol := portforward.PortForwardProtocol(wspf.PortForwardProtocolUDP)
	portforwardNew := &wspf.PortForwardNew{
		Protocol:   &protocol,
		RemoteHost: &p.remoteHost,
		RemotePort: &p.remotePort,
	}
	body, err := msgpack.Marshal(portforwardNew)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err.Error())
		panic(err)
	}
	m := &ws.ProtoMsg{
		Header: ws.ProtoHdr{
			Proto:     ws.ProtoTypePortForward,
			MsgType:   wspf.MessageTypePortForwardNew,
			SessionID: sessionID,
			Properties: map[string]interface{}{
				wspf.PropertyConnectionID: connectionID,
			},
		},
		Body: body,
	}
	msgChan <- m

	sendStopMessage := true
	defer func() {
		if sendStopMessage {
			m := &ws.ProtoMsg{
				Header: ws.ProtoHdr{
					Proto:     ws.ProtoTypePortForward,
					MsgType:   wspf.MessageTypePortForwardStop,
					SessionID: sessionID,
					Properties: map[string]interface{}{
						wspf.PropertyConnectionID: connectionID,
					},
				},
			}
			msgChan <- m
		}
	}()

	errChan := make(chan error)
	dataChan := make(chan []byte)

	// go routine to handle the network connection
	go p.handleRequestConnection(dataChan, errChan)

	// go routine to handle received messages
	go func() {
		for {
			select {
			case m := <-recvChan:
				if m.Header.Proto == ws.ProtoTypePortForward && m.Header.MsgType == wspf.MessageTypePortForwardStop {
					sendStopMessage = false
					return
				} else if m.Header.Proto == ws.ProtoTypePortForward && m.Header.MsgType == wspf.MessageTypePortForward {
					_, err := p.conn.WriteToUDP(m.Body, p.sourceAddr)
					if err != nil {
						fmt.Fprintf(os.Stderr, "error: %v\n", err.Error())
					} else {
						// send the ack
						m := &ws.ProtoMsg{
							Header: ws.ProtoHdr{
								Proto:     ws.ProtoTypePortForward,
								MsgType:   wspf.MessageTypePortForwardAck,
								SessionID: sessionID,
								Properties: map[string]interface{}{
									wspf.PropertyConnectionID: connectionID,
								},
							},
						}
						msgChan <- m
					}
				} else if m.Header.Proto == ws.ProtoTypePortForward && m.Header.MsgType == wspf.MessageTypePortForwardAck {
					p.waitGroupAcks.Add(-1)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// go routine to handle sent messages
	for {
		select {
		case err := <-errChan:
			if err != io.EOF {
				fmt.Fprintf(os.Stderr, "error: %v\n", err.Error())
			}
			return
		case data := <-dataChan:
			// wait to receive all the previous acks
			p.waitGroupAcks.Wait()

			// add an expected ack to the wait group
			p.waitGroupAcks.Add(1)

			m := &ws.ProtoMsg{
				Header: ws.ProtoHdr{
					Proto:     ws.ProtoTypePortForward,
					MsgType:   wspf.MessageTypePortForward,
					SessionID: sessionID,
					Properties: map[string]interface{}{
						wspf.PropertyConnectionID: connectionID,
					},
				},
				Body: data,
			}
			msgChan <- m
		case <-ctx.Done():
			return
		}
	}
}

func (p *UDPPortForwarder) handleRequestConnection(dataChan chan []byte, errChan chan error) {
	data := make([]byte, readBuffLength)

	for {
		n, udpAddr, err := p.conn.ReadFromUDP(data[:])
		if err != nil {
			errChan <- err
			break
		}
		if n > 0 {
			p.sourceAddr = udpAddr
			tmp := make([]byte, n)
			copy(tmp, data[:n])
			dataChan <- tmp
		}
	}
}
