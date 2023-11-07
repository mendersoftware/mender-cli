// Copyright 2023 Northern.tech AS
//
//	Licensed under the Apache License, Version 2.0 (the "License");
//	you may not use this file except in compliance with the License.
//	You may obtain a copy of the License at
//
//	    http://www.apache.org/licenses/LICENSE-2.0
//
//	Unless required by applicable law or agreed to in writing, software
//	distributed under the License is distributed on an "AS IS" BASIS,
//	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	See the License for the specific language governing permissions and
//	limitations under the License.
package cmd

import (
	"bufio"
	"compress/gzip"
	"context"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sys/unix"
	"golang.org/x/term"

	"github.com/mendersoftware/go-lib-micro/ws"
	wsshell "github.com/mendersoftware/go-lib-micro/ws/shell"

	"github.com/mendersoftware/mender-cli/client/deviceconnect"
	"github.com/mendersoftware/mender-cli/log"
)

const (
	// default terminal size
	defaultTermWidth  = 80
	defaultTermHeight = 40

	// dummy delay for playback
	playbackSleep = time.Millisecond * 32

	// cli args
	argRecord   = "record"
	argPlayback = "playback"
)

var terminalCmd = &cobra.Command{
	Use:   "terminal [DEVICE_ID]",
	Short: "Remotely access a terminal on a device",
	Long: "Remotely access a terminal on a device\n" +
		"Basic usage is terminal DEVICE_ID, which starts a new terminal " +
		"session with the remote device. The session can be saved locally " +
		"using --record flag. When using --playback flag, no DEVICE_ID is " +
		"required and no connection will be established.",
	Args: cobra.RangeArgs(0, 1),
	Run: func(c *cobra.Command, args []string) {
		cmd, err := NewTerminalCmd(c, args)
		CheckErr(err)
		CheckErr(cmd.Run())
	},
}

func init() {
	terminalCmd.Flags().StringP(argRecord, "", "", "recording file path to save the session to")
	terminalCmd.Flags().
		StringP(argPlayback, "", "", "recording file path to playback the session from")
}

// TerminalCmd handles the terminal command
type TerminalCmd struct {
	server             string
	token              string
	skipVerify         bool
	deviceID           string
	sessionID          string
	running            bool
	healthcheck        chan int
	stop               chan struct{}
	err                error
	recordFile         string
	recording          bool
	stopRecording      chan bool
	playbackFile       string
	terminalOutputChan chan []byte
}

const (
	deviceIDMaxLength     = 64
	terminalTypeMaxLength = 32
	terminalTypeDefault   = "xterm-256color"
)

type TerminalRecordingHeader struct {
	Version        uint8
	DeviceID       [deviceIDMaxLength]byte
	TerminalType   [terminalTypeMaxLength]byte
	TerminalWidth  int16
	TerminalHeight int16
	Timestamp      int64
}

const (
	terminalRecordingVersion = 1
)

type TerminalRecordingType int8

type TerminalRecordingData struct {
	Type TerminalRecordingType
	Data []byte
}

const (
	terminalRecordingOutput TerminalRecordingType = iota
)

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

	recordFile, err := cmd.Flags().GetString(argRecord)
	if err != nil {
		return nil, err
	}

	playbackFile, err := cmd.Flags().GetString(argPlayback)
	if err != nil {
		return nil, err
	}

	token, err := getAuthToken(cmd)
	if err != nil {
		return nil, err
	}

	deviceID := ""
	if len(args) == 1 {
		deviceID = args[0]
	}

	if playbackFile == "" && deviceID == "" {
		return nil, errors.New("No device specified")
	}

	return &TerminalCmd{
		server:             server,
		token:              token,
		skipVerify:         skipVerify,
		deviceID:           deviceID,
		healthcheck:        make(chan int),
		stop:               make(chan struct{}),
		recordFile:         recordFile,
		stopRecording:      make(chan bool),
		terminalOutputChan: make(chan []byte),
		playbackFile:       playbackFile,
	}, nil
}

// send the shell start message
func (c *TerminalCmd) startShell(client *deviceconnect.Client, termWidth, termHeight int) error {
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
	if err := client.WriteMessage(m); err != nil {
		return err
	}
	return nil
}

// send the stop shell message
func (c *TerminalCmd) stopShell(client *deviceconnect.Client) error {
	m := &ws.ProtoMsg{
		Header: ws.ProtoHdr{
			Proto:     ws.ProtoTypeShell,
			MsgType:   wsshell.MessageTypeStopShell,
			SessionID: c.sessionID,
		},
	}
	if err := client.WriteMessage(m); err != nil {
		return err
	}
	return nil
}

func (c *TerminalCmd) record() {
	f, err := os.Create(c.recordFile)
	if err != nil {
		log.Err(fmt.Sprintf("Can't create recording file: %s: %s", c.recordFile, err.Error()))
	}
	defer f.Close()

	fz := gzip.NewWriter(f)
	defer fz.Close()

	data := TerminalRecordingHeader{
		Version:        terminalRecordingVersion,
		Timestamp:      time.Now().Unix(),
		TerminalWidth:  defaultTermWidth,
		TerminalHeight: defaultTermHeight,
	}
	copy(data.DeviceID[:], []byte(c.deviceID))
	copy(data.TerminalType[:], []byte(terminalTypeDefault))
	err = binary.Write(fz, binary.LittleEndian, data)
	if err != nil {
		log.Err(fmt.Sprintf("Header write failed: %s", err.Error()))
	}
	err = fz.Flush()
	if err != nil {
		log.Err(fmt.Sprintf("Header flush failed: %s", err.Error()))
	}

	log.Info(fmt.Sprintf("Recording to file: %s", c.recordFile))

	e := gob.NewEncoder(fz)
	for {
		select {
		case <-c.stopRecording:
			return
		case terminalOutput := <-c.terminalOutputChan:
			o := TerminalRecordingData{
				Type: terminalRecordingOutput,
				Data: terminalOutput,
			}
			err = e.Encode(o)
			fz.Flush()
			if err != nil {
				log.Err(fmt.Sprintf("Error encoding %q: %s", string(terminalOutput), err.Error()))
				return
			}
		}
	}
}

func (c *TerminalCmd) playback(w io.Writer) error {
	f, err := os.Open(c.playbackFile)
	if err != nil {
		log.Err(fmt.Sprintf("Can't open %s: %s", c.playbackFile, err.Error()))
		return err
	}
	defer f.Close()

	fz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer fz.Close()

	var header TerminalRecordingHeader
	err = binary.Read(fz, binary.LittleEndian, &header)
	if err != nil {
		log.Err(fmt.Sprintf("Can't read header: %s", err.Error()))
		return err
	}

	dateTime := time.Unix(header.Timestamp, 0)

	log.Info(fmt.Sprintf("Playing back from file: %s", c.playbackFile))
	log.Info(fmt.Sprintf("Device ID: %s", string(header.DeviceID[:])))
	log.Info(fmt.Sprintf("Terminal type: %s", string(header.TerminalType[:])))
	log.Info(fmt.Sprintf("Terminal size: %dx%d", header.TerminalWidth, header.TerminalHeight))
	log.Info(fmt.Sprintf("Timestamp: %s", dateTime.Format(time.UnixDate)))
	log.Info("")

	d := gob.NewDecoder(fz)
	for {
		var o TerminalRecordingData
		err = d.Decode(&o)
		if err != nil {
			if err != io.EOF {
				log.Err(fmt.Sprintf("Decoding error: %s", err.Error()))
				return err
			}
			break
		}
		if o.Type == terminalRecordingOutput {
			_, err = w.Write(o.Data)
			if err != nil {
				log.Err(fmt.Sprintf("Writting error: %s", err.Error()))
				return err
			}
		}
		time.Sleep(playbackSleep)
	}
	log.Info("\r")
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

	// when playing back, no further processing is required
	if c.playbackFile != "" {
		if _, err := os.Stat(c.playbackFile); err == nil {
			return c.playback(os.Stdout)
		} else {
			return err
		}
	}

	// start recording when applicable
	if _, err := os.Stat(c.recordFile); os.IsNotExist(err) {
		if len(c.recordFile) > 0 {
			c.recording = true
			go c.record()
		}
	} else {
		log.Err(fmt.Sprintf(
			"Can't create recording file: %s exists, refused to record.",
			c.recordFile,
		))
	}

	client := deviceconnect.NewClient(c.server, c.token, c.skipVerify)

	// check if the device is connected
	device, err := client.GetDevice(c.deviceID)
	if err != nil {
		return errors.Wrap(err, "unable to get the device")
	} else if device.Status != deviceconnect.CONNECTED {
		return errors.New("the device is not connected")
	}

	// connect to the websocket and start the ping-pong connection health-check
	err = client.Connect(c.deviceID, c.token)
	if err != nil {
		return err
	}

	go client.PingPong(ctx)
	defer client.Close()

	// set the terminal in raw mode
	if term.IsTerminal(termID) {
		termWidth, termHeight, err = term.GetSize(termID)
		if err != nil {
			return errors.Wrap(err, "Unable to get the terminal size")
		}

		fmt.Fprintln(os.Stderr, "Press CTRL+] to quit the session")

		oldState, err := term.MakeRaw(termID)
		if err != nil {
			return errors.Wrap(err, "Unable to set the terminal in raw mode")
		}
		defer func() {
			_ = term.Restore(termID, oldState)
		}()
	}

	// start the shell
	if err := c.startShell(client, termWidth, termHeight); err != nil {
		return err
	}

	// wait for CTRL+C, signals or stop
	c.runLoop(ctx, client, termID, termWidth, termHeight)

	// cancel the context
	cancelContext()

	// stop shell message
	if err := c.stopShell(client); err != nil {
		return err
	}

	// return the error message (if any)
	return c.err
}

// Run executes the command
func (c *TerminalCmd) runLoop(
	ctx context.Context,
	client *deviceconnect.Client,
	termID, termWidth, termHeight int,
) {
	// message channel
	msgChan := make(chan *ws.ProtoMsg)

	c.running = true
	go c.pipeStdin(msgChan, os.Stdin)
	go c.pipeStdout(msgChan, client, os.Stdout)

	// handle CTRL+C and signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, unix.SIGINT, unix.SIGTERM)

	// resize the terminal window
	go c.resizeTerminal(ctx, msgChan, termID, termWidth, termHeight)

	healthcheckTimeout := time.Now().Add(24 * time.Hour)
	for c.running {
		select {
		case msg := <-msgChan:
			err := client.WriteMessage(msg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				break
			}
		case healthcheckInterval := <-c.healthcheck:
			healthcheckTimeout = time.Now().Add(time.Duration(healthcheckInterval) * time.Second)
		case <-time.After(time.Until(healthcheckTimeout)):
			_ = c.stopShell(client)
			c.err = errors.New("health check failed, connection with the device lost")
			c.running = false
		case <-quit:
			c.running = false
		case <-c.stop:
			c.running = false
		}
	}
}

func (c *TerminalCmd) resizeTerminal(
	ctx context.Context,
	msgChan chan *ws.ProtoMsg,
	termID int,
	termWidth int,
	termHeight int,
) {
	resize := make(chan os.Signal, 1)
	signal.Notify(resize, syscall.SIGWINCH)
	defer signal.Stop(resize)

	for {
		select {
		case <-ctx.Done():
			return
		case <-resize:
			newTermWidth, newTermHeight, _ := term.GetSize(termID)
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
	if c.recording {
		c.stopRecording <- true
	}
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

func (c *TerminalCmd) pipeStdout(
	msgChan chan *ws.ProtoMsg,
	client *deviceconnect.Client,
	w io.Writer,
) {
	for c.running {
		m, err := client.ReadMessage()
		if err != nil {
			if c.running {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
			} else {
				c.Stop()
			}
			break
		}
		if m.Header.Proto == ws.ProtoTypeShell &&
			m.Header.MsgType == wsshell.MessageTypeShellCommand {
			if _, err := w.Write(m.Body); err != nil {
				break
			}
			if c.recording {
				c.terminalOutputChan <- m.Body
			}
		} else if m.Header.Proto == ws.ProtoTypeShell &&
			m.Header.MsgType == wsshell.MessageTypePingShell {
			if healthcheckTimeout, ok := m.Header.Properties["timeout"].(int64); ok &&
				healthcheckTimeout > 0 {
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
		} else if m.Header.Proto == ws.ProtoTypeShell &&
			m.Header.MsgType == wsshell.MessageTypeSpawnShell {
			status, ok := m.Header.Properties["status"].(int64)
			if ok && status == int64(wsshell.ErrorMessage) {
				c.err = errors.New(fmt.Sprintf("Unable to start the shell: %s", string(m.Body)))
				c.Stop()
			} else {
				c.sessionID = string(m.Header.SessionID)
			}
		} else if m.Header.Proto == ws.ProtoTypeShell &&
			m.Header.MsgType == wsshell.MessageTypeStopShell {
			c.Stop()
			break
		}
	}
}
