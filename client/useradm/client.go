// Copyright 2018 Northern.tech AS
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
package useradm

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/mendersoftware/mender-cli/log"
)

const (
	loginUrl = "/api/management/v1/useradm/auth/login"
	timeout  = 10 * time.Second
)

type Client struct {
	url      string
	loginUrl string
	client   *http.Client
}

func NewClient(url string, skipVerify bool) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipVerify},
	}

	return &Client{
		url:      url,
		loginUrl: JoinURL(url, loginUrl),
		client: &http.Client{
			Transport: tr,
		},
	}
}

func (c *Client) Login(user, pass string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, c.loginUrl, nil)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req.SetBasicAuth(user, pass)

	reqDump, _ := httputil.DumpRequest(req, true)
	log.Verbf("sending request: \n%v", string(reqDump))

	rsp, err := c.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, errors.Wrap(err, "POST /auth/login request failed")
	}
	defer rsp.Body.Close()

	rspDump, _ := httputil.DumpResponse(rsp, true)
	log.Verbf("response: \n%v\n", string(rspDump))

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "can't read request body")
	}

	if rsp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("login failed with status %d", rsp.StatusCode))
	}

	return body, nil
}

func JoinURL(base, url string) string {
	if strings.HasPrefix(url, "/") {
		url = url[1:]
	}
	if !strings.HasSuffix(base, "/") {
		base = base + "/"
	}
	return base + url
}
