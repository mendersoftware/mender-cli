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
package deviceconnect

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mendersoftware/go-lib-micro/ws"
	"github.com/pkg/errors"
	"github.com/vmihailenco/msgpack"
)

const (
	// protocols
	httpProtocol = "http"
	wsProtocol   = "ws"

	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 1 * time.Minute
)

type Client struct {
	url        string
	skipVerify bool
	conn       *websocket.Conn
}

func NewClient(url string, skipVerify bool) *Client {
	return &Client{
		url:        url,
		skipVerify: skipVerify,
	}
}

// Connect to the websocket
func (c *Client) Connect(deviceID string, token []byte) error {
	fmt.Printf("Connecting to the device %s...\n", deviceID)
	deviceConnectPath := "/api/management/v1/deviceconnect/devices/" + deviceID + "/connect"
	u, err := url.Parse(strings.TrimSuffix(c.url, "/") + deviceConnectPath)
	if err != nil {
		return errors.Wrap(err, "Unable to parse the server URL")
	}
	u.Scheme = strings.Replace(u.Scheme, httpProtocol, wsProtocol, 1)

	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+string(token))
	websocket.DefaultDialer.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: c.skipVerify,
	}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), headers)
	if err != nil {
		return errors.Wrap(err, "Unable to connect to the device")
	}

	err = conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		return errors.Wrap(err, "Unable to set the read deadline")
	}

	c.conn = conn
	return nil
}

// PingPong handles the ping-pong connection health check
func (c *Client) PingPong(ctx context.Context) {
	pingPeriod := (pongWait * 9) / 10
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	c.conn.SetPongHandler(func(string) error {
		ticker.Reset(pingPeriod)
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	c.conn.SetPingHandler(func(msg string) error {
		ticker.Reset(pingPeriod)
		err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			return err
		}
		return c.conn.WriteControl(
			websocket.PongMessage,
			[]byte(msg),
			time.Now().Add(writeWait),
		)
	})

	for {
		select {
		case <-ticker.C:
			pongWaitString := strconv.Itoa(int(pongWait.Seconds()))
			_ = c.conn.WriteControl(
				websocket.PingMessage,
				[]byte(pongWaitString),
				time.Now().Add(writeWait),
			)

		case <-ctx.Done():
			return
		}
	}
}

// ReadMessage reads a Proto message from the websocket
func (c *Client) ReadMessage() (*ws.ProtoMsg, error) {
	_, data, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	m := &ws.ProtoMsg{}
	err = msgpack.Unmarshal(data, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// WriteMessage writes a Proto message to the websocket
func (c *Client) WriteMessage(m *ws.ProtoMsg) error {
	data, err := msgpack.Marshal(m)
	if err != nil {
		return errors.Wrap(err, "Unable to marshal the message from the websocket")
	}
	if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return errors.Wrap(err, "Unable to set the write deadline")
	}
	if err := c.conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		return errors.Wrap(err, "Unable to write the message")
	}
	return nil
}

// Close closes the connection
func (c *Client) Close() {
	c.conn.Close()
}
