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
package devices

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"

	"github.com/mendersoftware/mender-cli/client"
	"github.com/mendersoftware/mender-cli/log"
)

type deviceData struct {
	ID           string `json:"id"`
	IdentityData struct {
		Mac string `json:"mac"`
		Sku string `json:"sku"`
		Sn  string `json:"sn"`
	} `json:"identity_data"`
	Status    string `json:"status"`
	CreatedTs string `json:"created_ts"`
	UpdatedTs string `json:"updated_ts"`
	AuthSets  []struct {
		ID           string `json:"id"`
		PubKey       string `json:"pubkey"`
		IdentityData struct {
			Mac string `json:"mac"`
			Sku string `json:"sku"`
			Sn  string `json:"sn"`
		} `json:"identity_data"`
		Status string `json:"status"`
		Ts     string `json:"ts"`
	} `json:"auth_sets"`
	Decommissioning bool `json:"decommissioning"`
}

const (
	devicesListURL = "/api/management/v2/devauth/devices"
)

type Client struct {
	url            string
	devicesListURL string
	client         *http.Client
	output         io.Writer
}

func NewClient(url string, skipVerify bool) *Client {
	return &Client{
		url:            url,
		devicesListURL: client.JoinURL(url, devicesListURL),
		client:         client.NewHttpClient(skipVerify),
		output:         os.Stdout,
	}
}

func (c *Client) ListDevices(token string, detailLevel, perPage, page int, raw bool) error {
	if detailLevel > 3 || detailLevel < 0 {
		return fmt.Errorf("FAILURE: invalid devices detail")
	}

	req, err := http.NewRequest(http.MethodGet, c.devicesListURL, nil)
	if err != nil {
		return fmt.Errorf("failed to prepare request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	q := url.Values{
		"per_page": []string{strconv.Itoa(perPage)},
		"page":     []string{strconv.Itoa(page)},
	}
	req.URL.RawQuery = q.Encode()

	reqDump, err := httputil.DumpRequest(req, false)
	if err != nil {
		return err
	}
	log.Verbf("sending request: \n%s", string(reqDump))

	rsp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != 200 {
		return fmt.Errorf("GET %s request failed with status %d",
			req.URL.RequestURI(), rsp.StatusCode)
	}

	if raw {
		_, err := io.Copy(os.Stdout, rsp.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %w", err)
		}
	} else {
		var list []deviceData
		err = json.NewDecoder(rsp.Body).Decode(&list)
		if err != nil {
			return err
		}
		for _, v := range list {
			listDevice(c.output, v, detailLevel)
		}
	}
	return nil
}

func listDevice(out io.Writer, a deviceData, detailLevel int) {
	fmt.Fprintf(out, "ID: %s\n", a.ID)
	fmt.Fprintf(out, "Status: %s\n", a.Status)
	if detailLevel >= 1 {
		fmt.Println("IdentityData:")
		if a.IdentityData.Mac != "" {
			fmt.Fprintf(out, "  MAC address: %s\n", a.IdentityData.Mac)
		}
		if a.IdentityData.Sku != "" {
			fmt.Fprintf(out, "  Stock keeping unit: %s\n", a.IdentityData.Sku)
		}
		if a.IdentityData.Sn != "" {
			fmt.Fprintf(out, "  Serial number: %s\n", a.IdentityData.Sn)
		}
	}
	if detailLevel >= 1 {
		fmt.Fprintf(out, "CreatedTs: %s\n", a.CreatedTs)
		fmt.Fprintf(out, "UpdatedTs: %s\n", a.UpdatedTs)
		fmt.Fprintf(out, "Decommissioning: %t\n", a.Decommissioning)
	}
	if detailLevel >= 2 {
		for i, v := range a.AuthSets {
			fmt.Fprintf(out, "AuthSet[%d]:\n", i)
			fmt.Fprintf(out, "  ID: %s\n", v.ID)
			fmt.Fprintf(out, "  PubKey:\n%s", v.PubKey)
			fmt.Println("  IdentityData:")
			if v.IdentityData.Mac != "" {
				fmt.Fprintf(out, "    MAC address: %s\n", v.IdentityData.Mac)
			}
			if v.IdentityData.Sku != "" {
				fmt.Fprintf(out, "    Stock keeping unit: %s\n", v.IdentityData.Sku)
			}
			if v.IdentityData.Sn != "" {
				fmt.Fprintf(out, "    Serial number: %s\n", v.IdentityData.Sn)
			}
			fmt.Fprintf(out, "  Status: %s\n", v.Status)
			fmt.Fprintf(out, "  Ts: %s\n", v.Ts)
		}
	}

	fmt.Fprintf(
		out, "--------------------------------------------------------------------------------",
	)
}
