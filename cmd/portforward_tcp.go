// Copyright 2022 Northern.tech AS
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
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/mendersoftware/go-lib-micro/ws"
	"github.com/mendersoftware/go-lib-micro/ws/portforward"
	wspf "github.com/mendersoftware/go-lib-micro/ws/portforward"
	"github.com/vmihailenco/msgpack"
)

const portForwardTCPChannelSize = 20

type TCPPortForwarder struct {
	listen     net.Listener
	remoteHost string
	remotePort uint16
}

func NewTCPPortForwarder(
	bindingHost string,
	localPort uint16,
	remoteHost string,
	remotePort uint16,
) (*TCPPortForwarder, error) {
	fmt.Printf("Forwarding from %s:%d -> %s:%d\n", bindingHost, localPort, remoteHost, remotePort)
	listen, err := net.Listen(protocolTCP, bindingHost+":"+strconv.Itoa(int(localPort)))
	if err != nil {
		return nil, err
	}
	return &TCPPortForwarder{
		listen:     listen,
		remoteHost: remoteHost,
		remotePort: remotePort,
	}, nil
}

func (p *TCPPortForwarder) Run(
	ctx context.Context,
	sessionID string,
	msgChan chan *ws.ProtoMsg,
	recvChans map[string]chan *ws.ProtoMsg,
) {
	// listen for new connections
	defer p.listen.Close()
	acceptedConnections := make(chan net.Conn)

	// go-routine to accept new connections
	go func() {
		for {
			conn, err := p.listen.Accept()
			if err != nil {
				return
			}
			fmt.Printf(
				"Handling connection from %s to %s\n",
				conn.RemoteAddr().String(),
				conn.LocalAddr().String(),
			)
			acceptedConnections <- conn
		}
	}()

	// handle new connections
	for {
		select {
		case conn := <-acceptedConnections:
			connectionUUID, _ := uuid.NewUUID()
			connectionID := connectionUUID.String()
			recvChan := make(chan *ws.ProtoMsg, portForwardTCPChannelSize)
			recvChans[connectionID] = recvChan
			go p.handleRequest(ctx, conn, sessionID, connectionID, recvChan, msgChan)
		case <-ctx.Done():
			return
		}
	}
}

func (p *TCPPortForwarder) handleRequest(
	ctx context.Context,
	conn net.Conn,
	sessionID string,
	connectionID string,
	recvChan chan *ws.ProtoMsg,
	msgChan chan *ws.ProtoMsg,
) {
	defer conn.Close()

	ackChan := make(chan struct{})
	defer func() { close(ackChan) }()

	errChan := make(chan error)
	dataChan := make(chan []byte)

	protocol := portforward.PortForwardProtocol(wspf.PortForwardProtocolTCP)
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
		conn.Close()
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

	// go routine to handle the network connection
	go p.handleRequestConnection(dataChan, errChan, conn)

	// go routine to handle received messages
	go func(connectionID string) {
		for {
			select {
			case m := <-recvChan:
				if m.Header.Proto == ws.ProtoTypePortForward &&
					m.Header.MsgType == wspf.MessageTypePortForwardStop {
					sendStopMessage = false
					return
				} else if m.Header.Proto == ws.ProtoTypePortForward &&
					m.Header.MsgType == wspf.MessageTypePortForward {
					_, err := conn.Write(m.Body)
					if err != nil {
						if errors.Unwrap(err) != net.ErrClosed {
							fmt.Fprintf(os.Stderr, "error: %v\n", err.Error())
						}
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
				} else if m.Header.Proto == ws.ProtoTypePortForward &&
					m.Header.MsgType == wspf.MessageTypePortForwardAck {
					<-ackChan
				}
			case <-ctx.Done():
				return
			}
		}
	}(connectionID)

	// go routine to handle sent messages
	for err == nil {
		select {
		case err = <-errChan:

		case data := <-dataChan:
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
			// wait for the ack to be received before processing more data
			select {
			case ackChan <- struct{}{}:
			case <-ctx.Done():
				err = ctx.Err()

			case err = <-errChan:

			}
		case <-ctx.Done():
			err = ctx.Err()
		}
	}

	if err != io.EOF {
		fmt.Fprintf(os.Stderr, "error: %v\n", err.Error())
	}
}

func (p *TCPPortForwarder) handleRequestConnection(
	dataChan chan []byte,
	errChan chan error,
	conn net.Conn,
) {
	data := make([]byte, readBuffLength)

	for {
		n, err := conn.Read(data)
		if err != nil {
			errChan <- err
			break
		}
		if n > 0 {
			tmp := make([]byte, n)
			copy(tmp, data[:n])
			dataChan <- tmp
		}
	}
}
