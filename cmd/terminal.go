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
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vmihailenco/msgpack"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/sys/unix"
	"golang.org/x/term"

	"github.com/mendersoftware/go-lib-micro/ws"
	wsshell "github.com/mendersoftware/go-lib-micro/ws/shell"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 1 * time.Minute

	// protocols
	httpProtocol = "http"
	wsProtocol   = "ws"

	// default terminal size
	defaultTermWidth  = 80
	defaultTermHeight = 40
)

var terminalCmd = &cobra.Command{
	Use:   "terminal DEVICE_ID",
	Short: "Access a device's remote terminal",
	Args:  cobra.ExactArgs(1),
	Run: func(c *cobra.Command, args []string) {
		cmd, err := NewTerminalCmd(c, args)
		CheckErr(err)
		CheckErr(cmd.Run())
	},
}

// TerminalCmd handles the terminal command
type TerminalCmd struct {
	server      string
	skipVerify  bool
	deviceID    string
	sessionID   string
	running     bool
	healthcheck chan int
	stop        chan struct{}
	err         error
}

// NewTerminalCmd returns a new TerminalCmd
func NewTerminalCmd(cmd *cobra.Command, args []string) (*TerminalCmd, error) {
	server := viper.GetString(argRootServer)
	if server == "" {
		return nil, errors.New("No server")
	}

	skipVerify, err := cmd.Flags().GetBool(argRootSkipVerify)
	if err != nil {
		return nil, err
	}

	return &TerminalCmd{
		server:      server,
		skipVerify:  skipVerify,
		deviceID:    args[0],
		healthcheck: make(chan int),
		stop:        make(chan struct{}),
	}, nil
}

func (c *TerminalCmd) getToken() ([]byte, error) {
	tokenPath, err := getDefaultAuthTokenPath()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to determine the auth token path")
	}
	token, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return nil, errors.Wrap(err, "Please Login first")
	}
	return token, nil
}

// connect to the websocket
func (c *TerminalCmd) connect() (*websocket.Conn, error) {
	token, err := c.getToken()
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(os.Stderr, "Connecting to the remote terminal of the device %s...\n", c.deviceID)
	deviceConnectPath := "/api/management/v1/deviceconnect/devices/" + c.deviceID + "/connect"
	u, err := url.Parse(strings.TrimSuffix(c.server, "/") + deviceConnectPath)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to parse the server URL")
	}
	u.Scheme = strings.Replace(u.Scheme, httpProtocol, wsProtocol, 1)

	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+string(token))
	websocket.DefaultDialer.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: c.skipVerify,
	}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), headers)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to connect to the device")
	}

	err = conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		return nil, errors.Wrap(err, "Unable to set the read deadline")
	}

	return conn, nil
}

// handle the ping-pong connection health check
func (c *TerminalCmd) pingPong(ctx context.Context, conn *websocket.Conn) {
	pingPeriod := (pongWait * 9) / 10
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	conn.SetPongHandler(func(string) error {
		ticker.Reset(pingPeriod)
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	conn.SetPingHandler(func(msg string) error {
		ticker.Reset(pingPeriod)
		err := conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			return err
		}
		return conn.WriteControl(
			websocket.PongMessage,
			[]byte(msg),
			time.Now().Add(writeWait),
		)
	})

	for {
		select {
		case <-ticker.C:
			pongWaitString := strconv.Itoa(int(pongWait.Seconds()))
			_ = conn.WriteControl(
				websocket.PingMessage,
				[]byte(pongWaitString),
				time.Now().Add(writeWait),
			)

		case <-ctx.Done():
			return
		}
	}
}

// send the shell start message
func (c *TerminalCmd) startShell(conn *websocket.Conn, termWidth, termHeight int) error {
	m := &ws.ProtoMsg{
		Header: ws.ProtoHdr{
			Proto:   ws.ProtoTypeShell,
			MsgType: wsshell.MessageTypeSpawnShell,
			Properties: map[string]interface{}{
				"terminal_width":  termWidth,
				"terminal_height": termHeight,
			},
		},
	}

	data, err := msgpack.Marshal(m)
	if err != nil {
		return errors.Wrap(err, "Unable to marshal the message from the websocket")
	}
	if err := conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return errors.Wrap(err, "Unable to set the write deadline")
	}
	if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		return errors.Wrap(err, "Unable to write the message")
	}
	return nil
}

// send the stop shell message
func (c *TerminalCmd) stopShell(conn *websocket.Conn) error {
	m := &ws.ProtoMsg{
		Header: ws.ProtoHdr{
			Proto:     ws.ProtoTypeShell,
			MsgType:   wsshell.MessageTypeStopShell,
			SessionID: c.sessionID,
		},
	}

	data, err := msgpack.Marshal(m)
	if err != nil {
		return errors.Wrap(err, "Unable to marshal the message from the websocket")
	}
	if err := conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return errors.Wrap(err, "Unable to set the write deadline")
	}
	if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		return errors.Wrap(err, "Unable to write the message")
	}
	return nil
}

// Run executes the command
func (c *TerminalCmd) Run() error {
	ctx, cancelContext := context.WithCancel(context.Background())
	defer cancelContext()

	// get the terminal width and height
	termWidth := defaultTermWidth
	termHeight := defaultTermHeight
	termID := int(os.Stdout.Fd())
	isTerminal := false

	stat, _ := os.Stdout.Stat()
	if (stat.Mode() & os.ModeCharDevice) > 0 {
		var err error
		termWidth, termHeight, err = terminal.GetSize(termID)
		if err != nil {
			return errors.Wrap(err, "Unable to get the terminal size")
		}
		isTerminal = true
	}

	// connect to the websocket
	conn, err := c.connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	// set the terminal in raw mode
	if isTerminal {
		fmt.Fprintln(os.Stderr, "Press CTRL+] to quit the session")

		oldState, err := term.MakeRaw(termID)
		if err != nil {
			return errors.Wrap(err, "Unable to set the terminal in raw mode")
		}
		defer func() {
			_ = term.Restore(termID, oldState)
		}()
	}

	// start the ping-pong connection health-check
	go c.pingPong(ctx, conn)

	// start the shell
	if err := c.startShell(conn, termWidth, termHeight); err != nil {
		return err
	}

	// message channel
	msgChan := make(chan *ws.ProtoMsg)

	c.running = true
	go c.pipeStdin(msgChan, os.Stdin)
	go c.pipeStdout(msgChan, conn, os.Stdout)

	// handle CTRL+C and signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, unix.SIGINT, unix.SIGTERM)

	// resize the terminal window
	go c.resizeTerminal(ctx, msgChan, termID, termWidth, termHeight)

	// wait for CTRL+C, signals or stop
	healthcheckTimeout := time.Now().Add(24 * time.Hour)
	for c.running {
		select {
		case msg := <-msgChan:
			data, err := msgpack.Marshal(msg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				break
			}
			_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
			_ = conn.WriteMessage(websocket.BinaryMessage, data)
		case healthcheckInterval := <-c.healthcheck:
			healthcheckTimeout = time.Now().Add(time.Duration(healthcheckInterval) * time.Second)
		case <-time.After(time.Until(healthcheckTimeout)):
			_ = c.stopShell(conn)
			c.err = errors.New("health check failed, connection with the device lost")
			c.running = false
		case <-quit:
			c.running = false
		case <-c.stop:
			c.running = false
		}
	}

	// cancel the context
	cancelContext()

	// stop shell message
	if err := c.stopShell(conn); err != nil {
		return err
	}

	// return the error message (if any)
	return c.err
}

func (c *TerminalCmd) resizeTerminal(ctx context.Context, msgChan chan *ws.ProtoMsg, termID int, termWidth int, termHeight int) {
	resize := make(chan os.Signal, 1)
	signal.Notify(resize, syscall.SIGWINCH)
	defer signal.Stop(resize)

	for {
		select {
		case <-ctx.Done():
			return
		case <-resize:
			newTermWidth, newTermHeight, _ := terminal.GetSize(termID)
			if newTermWidth != termWidth || newTermHeight != termHeight {
				termWidth = newTermWidth
				termHeight = newTermHeight
				m := &ws.ProtoMsg{
					Header: ws.ProtoHdr{
						Proto:   ws.ProtoTypeShell,
						MsgType: wsshell.MessageTypeResizeShell,
						Properties: map[string]interface{}{
							"terminal_width":  termWidth,
							"terminal_height": termHeight,
						},
					},
				}
				msgChan <- m
			}
		}
	}
}

func (c *TerminalCmd) Stop() {
	c.running = false
	c.stop <- struct{}{}
}

func (c *TerminalCmd) pipeStdin(msgChan chan *ws.ProtoMsg, r io.Reader) {
	s := bufio.NewReader(r)
	for c.running {
		raw := make([]byte, 1024)
		n, err := s.Read(raw)
		if err != nil {
			if c.running {
				if err != io.EOF {
					fmt.Fprintf(os.Stderr, "error: %v\n", err)
				}
			} else {
				c.Stop()
			}
			break
		}
		// CTRL+] terminates the session
		if raw[0] == 29 {
			c.Stop()
			return
		}

		m := &ws.ProtoMsg{
			Header: ws.ProtoHdr{
				Proto:     ws.ProtoTypeShell,
				MsgType:   wsshell.MessageTypeShellCommand,
				SessionID: c.sessionID,
			},
			Body: raw[:n],
		}
		msgChan <- m
	}
}

func (c *TerminalCmd) pipeStdout(msgChan chan *ws.ProtoMsg, conn *websocket.Conn, w io.Writer) {
	for c.running {
		_, data, err := conn.ReadMessage()
		if err != nil {
			if c.running {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
			} else {
				c.Stop()
			}
			break
		}

		m := &ws.ProtoMsg{}
		err = msgpack.Unmarshal(data, m)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			break
		}
		if m.Header.Proto == ws.ProtoTypeShell && m.Header.MsgType == wsshell.MessageTypeShellCommand {
			if _, err := w.Write(m.Body); err != nil {
				break
			}
		} else if m.Header.Proto == ws.ProtoTypeShell && m.Header.MsgType == wsshell.MessageTypePingShell {
			if healthcheckTimeout, ok := m.Header.Properties["timeout"].(int64); ok && healthcheckTimeout > 0 {
				c.healthcheck <- int(healthcheckTimeout)
			}
			m := &ws.ProtoMsg{
				Header: ws.ProtoHdr{
					Proto:     ws.ProtoTypeShell,
					MsgType:   wsshell.MessageTypePongShell,
					SessionID: c.sessionID,
				},
			}
			msgChan <- m
		} else if m.Header.Proto == ws.ProtoTypeShell && m.Header.MsgType == wsshell.MessageTypeSpawnShell {
			status, ok := m.Header.Properties["status"].(int64)
			if ok && status == int64(wsshell.ErrorMessage) {
				c.err = errors.New(fmt.Sprintf("Unable to start the shell: %s", string(m.Body)))
				c.Stop()
			} else {
				c.sessionID = string(m.Header.SessionID)
			}
		} else if m.Header.Proto == ws.ProtoTypeShell && m.Header.MsgType == wsshell.MessageTypeStopShell {
			c.Stop()
			break
		}
	}
}
