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
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/mendersoftware/go-lib-micro/ws"
	wspf "github.com/mendersoftware/go-lib-micro/ws/portforward"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/vmihailenco/msgpack"
	"golang.org/x/sys/unix"

	"github.com/mendersoftware/mender-cli/client/deviceconnect"
)

const (
	argBindHost    = "bind"
	readBuffLength = 4096
	localhost      = "127.0.0.1"
)

var portForwardCmd = &cobra.Command{
	Use: "port-forward DEVICE_ID [tcp|udp/]LOCAL_PORT[:REMOTE_PORT]" +
		" [[tcp|udp/]LOCAL_PORT[:REMOTE_PORT]...]",
	Short: "Forward one or more local ports to remote port(s) on the device",
	Long: "This command supports both TCP and UDP port-forwarding.\n\n" +
		"The port specification can be prefixed with \"tcp/\" or \"udp/\".\n" +
		"If no prefix is specified, TCP is the default.\n\n" +
		"REMOTE_PORT can also be specified in the form REMOTE_HOST:REMOTE_PORT, making\n" +
		"it possible to port-forward to third hosts running in the device's network.\n" +
		"In this case, the specification will be LOCAL_PORT:REMOTE_HOST:REMOTE_PORT.\n\n" +
		"You can specify multiple port mapping specifications.",
	Example: "  mender-cli port-forward DEVICE_ID 8000:8000\n" +
		"  mender-cli port-forward DEVICE_ID udp/8000:8000\n" +
		"  mender-cli port-forward DEVICE_ID tcp/8000:192.168.1.1:8000",
	Args: cobra.MinimumNArgs(2),
	Run: func(c *cobra.Command, args []string) {
		cmd, err := NewPortForwardCmd(c, args)
		CheckErr(err)
		CheckErr(cmd.Run())
	},
}

var portForwardMaxDuration = 24 * time.Hour

var errPortForwardNotImplemented = errors.New(
	"port forward not implemented or enabled on the device",
)
var errRestart = errors.New("restart")

func init() {
	portForwardCmd.Flags().StringP(argBindHost, "", localhost, "binding host")
}

const (
	protocolTCP = "tcp"
	protocolUDP = "udp"
)

type portMapping struct {
	Protocol   string
	LocalPort  uint16
	RemoteHost string
	RemotePort uint16
}

// PortForwardCmd handles the port-forward command
type PortForwardCmd struct {
	server       string
	skipVerify   bool
	deviceID     string
	sessionID    string
	bindingHost  string
	portMappings []portMapping
	recvChans    map[string]chan *ws.ProtoMsg
	running      bool
	stop         chan struct{}
	err          error
}

func getPortMappings(args []string) ([]portMapping, error) {
	var err error
	portMappings := []portMapping{}
	for _, arg := range args {
		remoteHost := localhost
		protocol := wspf.PortForwardProtocolTCP
		if strings.Contains(arg, "/") {
			parts := strings.SplitN(arg, "/", 2)
			if parts[0] == protocolTCP {
				protocol = protocolTCP
			} else if parts[0] == protocolUDP {
				protocol = protocolUDP
			} else {
				return nil, errors.New("unknown protocol: " + parts[0])
			}
			arg = parts[1]
		}
		var localPort, remotePort int
		if strings.Contains(arg, ":") {
			parts := strings.SplitN(arg, ":", 3)
			if len(parts) == 3 {
				remoteHost = parts[1]
				parts = []string{parts[0], parts[2]}
			}
			localPort, err = strconv.Atoi(parts[0])
			if err != nil || localPort < 0 || localPort > 65536 {
				return nil, errors.New("invalid port number: " + parts[0])
			}
			remotePort, err = strconv.Atoi(parts[1])
			if err != nil || remotePort < 0 || remotePort > 65536 {
				return nil, errors.New("invalid port number: " + parts[1])
			}
		} else {
			port, err := strconv.Atoi(arg)
			if err != nil || port < 0 || port > 65536 {
				return nil, errors.New("invalid port number: " + arg)
			}
			localPort = port
			remotePort = port
		}
		portMappings = append(portMappings, portMapping{
			Protocol:   protocol,
			LocalPort:  uint16(localPort),
			RemoteHost: remoteHost,
			RemotePort: uint16(remotePort),
		})
	}
	return portMappings, nil
}

// NewPortForwardCmd returns a new PortForwardCmd
func NewPortForwardCmd(cmd *cobra.Command, args []string) (*PortForwardCmd, error) {
	server, err := cmd.Flags().GetString(argRootServer)
	if err != nil {
		return nil, err
	}

	skipVerify, err := cmd.Flags().GetBool(argRootSkipVerify)
	if err != nil {
		return nil, err
	}

	bindingHost, err := cmd.Flags().GetString(argBindHost)
	if err != nil {
		return nil, err
	}

	portMappings, err := getPortMappings(args[1:])
	if err != nil {
		return nil, err
	}

	return &PortForwardCmd{
		server:       server,
		skipVerify:   skipVerify,
		deviceID:     args[0],
		bindingHost:  bindingHost,
		portMappings: portMappings,
		recvChans:    make(map[string]chan *ws.ProtoMsg),
		stop:         make(chan struct{}),
	}, nil
}

func (c *PortForwardCmd) getToken(tokenPath string) ([]byte, error) {
	token, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return nil, errors.Wrap(err, "Please Login first")
	}
	return token, nil
}

// Run executes the command
func (c *PortForwardCmd) Run() error {
	tokenPath, err := getDefaultAuthTokenPath()
	if err != nil {
		return errors.Wrap(err, "Unable to determine the default auth token path")
	}

	for {
		if err := c.run(tokenPath); err != errRestart {
			return err
		}
	}
}

func (c *PortForwardCmd) run(tokenPath string) error {
	ctx, cancelContext := context.WithCancel(context.Background())
	defer cancelContext()

	client := deviceconnect.NewClient(c.server, tokenPath, c.skipVerify)

	// check if the device is connected
	device, err := client.GetDevice(c.deviceID)
	if err != nil {
		return errors.Wrap(err, "unable to get the device")
	} else if device.Status != deviceconnect.CONNECTED {
		return errors.New("the device is not connected")
	}

	// get the JWT token
	token, err := c.getToken(tokenPath)
	if err != nil {
		return err
	}

	// connect to the websocket and start the ping-pong connection health-check
	err = client.Connect(c.deviceID, token)
	if err != nil {
		return err
	}

	go client.PingPong(ctx)
	defer client.Close()

	// perform ws protocol handshake
	err = c.handshake(client)
	if err != nil {
		return err
	}

	// message channel
	msgChan := make(chan *ws.ProtoMsg)

	// start the local TCP listeners
	for _, portMapping := range c.portMappings {
		switch portMapping.Protocol {
		case protocolTCP:
			forwarder, err := NewTCPPortForwarder(c.bindingHost, portMapping.LocalPort,
				portMapping.RemoteHost, portMapping.RemotePort)
			if err != nil {
				return err
			}
			go forwarder.Run(ctx, c.sessionID, msgChan, c.recvChans)
		case protocolUDP:
			forwarder, err := NewUDPPortForwarder(c.bindingHost, portMapping.LocalPort,
				portMapping.RemoteHost, portMapping.RemotePort)
			if err != nil {
				return err
			}
			go forwarder.Run(ctx, c.sessionID, msgChan, c.recvChans)
		default:
			return errors.New("unknown protocol: " + portMapping.Protocol)
		}
	}

	c.running = true
	go c.processIncomingMessages(msgChan, client)

	// handle CTRL+C and signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, unix.SIGINT, unix.SIGTERM)

	// wait for CTRL+C, signals or stop
	restart := false
	timeout := time.Now().Add(portForwardMaxDuration)
	for c.running {
		select {
		case msg := <-msgChan:
			err := client.WriteMessage(msg)
			if err != nil {
				c.err = err
				break
			}
		case <-time.After(time.Until(timeout)):
			c.err = errors.New("port forward timed out: max duration reached")
			c.running = false
		case <-quit:
			c.running = false
		case <-c.stop:
			restart = true
			c.running = false
		}
	}

	// cancel the context
	cancelContext()

	// close the ws session
	err = c.closeSession(client)
	if c.err == nil && err != nil {
		c.err = err
	}

	// if stopping because of an error, restart the port-forwarding command
	if restart {
		return errRestart
	}

	// return the error message (if any)
	return c.err
}

func (c *PortForwardCmd) Stop() {
	c.stop <- struct{}{}
}

// handshake initiates a handshake and checks that the device
// is willing to accept port forward requests.
func (c *PortForwardCmd) handshake(client *deviceconnect.Client) error {
	// open the session
	body, err := msgpack.Marshal(&ws.Open{
		Versions: []int{ws.ProtocolVersion},
	})
	if err != nil {
		return err
	}
	m := &ws.ProtoMsg{
		Header: ws.ProtoHdr{
			Proto:   ws.ProtoTypeControl,
			MsgType: ws.MessageTypeOpen,
		},
		Body: body,
	}
	err = client.WriteMessage(m)
	if err != nil {
		return err
	}

	msg, err := client.ReadMessage()
	if err != nil {
		return err
	}
	if msg.Header.MsgType == ws.MessageTypeError {
		erro := new(ws.Error)
		_ = msgpack.Unmarshal(msg.Body, erro)
		return errors.Errorf("handshake error from client: %s", erro.Error)
	} else if msg.Header.MsgType != ws.MessageTypeAccept {
		return errPortForwardNotImplemented
	}

	accept := new(ws.Accept)
	err = msgpack.Unmarshal(msg.Body, accept)
	if err != nil {
		return err
	}

	found := false
	for _, proto := range accept.Protocols {
		if proto == ws.ProtoTypePortForward {
			found = true
			break
		}
	}
	if !found {
		return errPortForwardNotImplemented
	}

	c.sessionID = msg.Header.SessionID
	return nil
}

// closeSession closes the WS session
func (c *PortForwardCmd) closeSession(client *deviceconnect.Client) error {
	m := &ws.ProtoMsg{
		Header: ws.ProtoHdr{
			Proto:   ws.ProtoTypeControl,
			MsgType: ws.MessageTypeClose,
		},
	}
	err := client.WriteMessage(m)
	if err != nil {
		return err
	}

	return nil
}

func (c *PortForwardCmd) processIncomingMessages(
	msgChan chan *ws.ProtoMsg,
	client *deviceconnect.Client,
) {
	for c.running {
		m, err := client.ReadMessage()
		if err != nil {
			c.err = err
			c.Stop()
			break
		} else if m.Header.Proto == ws.ProtoTypeControl && m.Header.MsgType == ws.MessageTypePing {
			m := &ws.ProtoMsg{
				Header: ws.ProtoHdr{
					Proto:     ws.ProtoTypeControl,
					MsgType:   ws.MessageTypePong,
					SessionID: c.sessionID,
				},
			}
			msgChan <- m
		} else if m.Header.Proto == ws.ProtoTypePortForward &&
			m.Header.MsgType == ws.MessageTypeError {
			erro := new(ws.Error)
			if err := msgpack.Unmarshal(m.Body, erro); err != nil &&
				erro.MessageType != wspf.MessageTypePortForwardStop {
				c.err = errors.New(fmt.Sprintf(
					"Unable to start the port-forwarding: %s",
					string(m.Body),
				))
				c.running = false
				c.Stop()
			}
		} else if m.Header.Proto == ws.ProtoTypePortForward &&
			(m.Header.MsgType == wspf.MessageTypePortForward ||
				m.Header.MsgType == wspf.MessageTypePortForwardAck ||
				m.Header.MsgType == wspf.MessageTypePortForwardStop) {
			connectionID, _ := m.Header.Properties[wspf.PropertyConnectionID].(string)
			if connectionID != "" {
				if recvChan, ok := c.recvChans[connectionID]; ok {
					recvChan <- m
				}
			}
		}
	}
}
