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
package useradm

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/pkg/errors"

	"github.com/mendersoftware/mender-cli/client"
	"github.com/mendersoftware/mender-cli/log"
)

const (
	loginUrl = "/api/management/v1/useradm/auth/login"
	meUrl    = "/api/management/v1/useradm/users/me"
	timeout  = 10 * time.Second
)

type Client struct {
	url      string
	loginUrl string
	client   *http.Client
}

func NewClient(url string, skipVerify bool) *Client {
	return &Client{
		url:      url,
		loginUrl: client.JoinURL(url, loginUrl),
		client:   client.NewHttpClient(skipVerify),
	}
}

func (c *Client) Login(user, pass string, token string) ([]byte, error) {
	var reader *bytes.Reader
	var req *http.Request
	var err error

	if len(token) > 1 {
		tokenJson := "{\"token2fa\":\"" + token + "\"}"
		reader = bytes.NewReader([]byte(tokenJson))
		req, err = http.NewRequest(http.MethodPost, c.loginUrl, reader)
	} else {
		req, err = http.NewRequest(http.MethodPost, c.loginUrl, nil)
	}
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req.SetBasicAuth(user, pass)
	req.Header.Set("Content-Type", "application/json")

	reqDump, _ := httputil.DumpRequest(req, true)
	log.Verbf("sending request: \n%v", string(reqDump))

	rsp, err := c.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, errors.Wrap(err, "POST /auth/login request failed")
	}
	defer rsp.Body.Close()

	rspDump, _ := httputil.DumpResponse(rsp, true)
	log.Verbf("response: \n%v\n", string(rspDump))

	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "can't read request body")
	}

	if rsp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("login failed with status %d", rsp.StatusCode))
	}

	return body, nil
}

// VerifyError is returned by Verify when the server explicitly rejects the
// token (HTTP 401/403), as opposed to a transport-level failure.
type VerifyError struct {
	StatusCode int
}

func (e *VerifyError) Error() string {
	return fmt.Sprintf("token rejected by server (status %d)", e.StatusCode)
}

// Verify checks that the given token is accepted by the configured server by
// calling GET /useradm/users/me. It returns nil on HTTP 200, a *VerifyError on
// 401/403, and a generic error for any other transport or status failure.
func (c *Client) Verify(token string) error {
	req, err := http.NewRequest(http.MethodGet, client.JoinURL(c.url, meUrl), nil)
	if err != nil {
		return errors.Wrap(err, "failed to create verify request")
	}
	req.Header.Set("Authorization", "Bearer "+token)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	reqDump, _ := httputil.DumpRequest(req, false)
	log.Verbf("sending request: \n%v", string(reqDump))

	rsp, err := c.client.Do(req.WithContext(ctx))
	if err != nil {
		return errors.Wrap(err, "GET /useradm/users/me request failed")
	}
	defer rsp.Body.Close()

	rspDump, _ := httputil.DumpResponse(rsp, true)
	log.Verbf("response: \n%v\n", string(rspDump))

	switch rsp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusUnauthorized, http.StatusForbidden:
		return &VerifyError{StatusCode: rsp.StatusCode}
	default:
		return fmt.Errorf("verify failed with status %d", rsp.StatusCode)
	}
}
