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
package client

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/pkg/errors"

	"github.com/mendersoftware/mender-cli/log"
)

func NewHttpClient(skipVerify bool) *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipVerify},
	}

	return &http.Client{
		Transport: tr,
	}
}

func JoinURL(base, url string) string {
	url = strings.TrimPrefix(url, "/")
	if !strings.HasSuffix(base, "/") {
		base = base + "/"
	}
	return base + url
}

func DoGetRequest(token, urlPath string, client *http.Client) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create HTTP request")
	}
	req.Header.Set("Authorization", "Bearer "+string(token))

	reqDump, err := httputil.DumpRequest(req, false)
	if err != nil {
		return nil, err
	}
	log.Verbf("sending request: \n%s", string(reqDump))

	rsp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Get %s request failed", urlPath))
	}
	if rsp.StatusCode != 200 {
		return nil, fmt.Errorf("Get %s request failed with status %d\n", urlPath, rsp.StatusCode)
	}

	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
