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
package devices

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/pkg/errors"

	"github.com/mendersoftware/mender-cli/log"
)

type devicesList struct {
	devices []deviceData
}

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
}

func NewClient(url string, skipVerify bool) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipVerify},
	}

	return &Client{
		url:            url,
		devicesListURL: JoinURL(url, devicesListURL),
		client: &http.Client{
			Transport: tr,
		},
	}
}

func (c *Client) ListDevices(tokenPath string, detailLevel int) error {
	if detailLevel > 3 || detailLevel < 0 {
		return fmt.Errorf("FAILURE: invalid devices detail")
	}
	token, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return errors.Wrap(err, "Please Login first")
	}

	req, err := http.NewRequest(http.MethodGet, c.devicesListURL, nil)
	if err != nil {
		return errors.Wrap(err, "Failed to create HTTP request")
	}
	req.Header.Set("Authorization", "Bearer "+string(token))

	reqDump, err := httputil.DumpRequest(req, false)
	if err != nil {
		return err
	}
	log.Verbf("sending request: \n%s", string(reqDump))

	rsp, err := c.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "Get /devauth/devices request failed")
	}
	if rsp.StatusCode != 200 {
		return fmt.Errorf("Get /devauth/devices request failed with status %d\n", rsp.StatusCode)
	}

	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	var list devicesList
	err = json.Unmarshal(body, &list.devices)
	if err != nil {
		return err
	}
	for _, v := range list.devices {
		listDevice(v, detailLevel)
	}

	return nil
}

func listDevice(a deviceData, detailLevel int) {
	fmt.Printf("ID: %s\n", a.ID)
	fmt.Printf("Status: %s\n", a.Status)
	if detailLevel >= 1 {
		fmt.Println("IdentityData:")
		if a.IdentityData.Mac != "" {
			fmt.Printf("  MAC address: %s\n", a.IdentityData.Mac)
		}
		if a.IdentityData.Sku != "" {
			fmt.Printf("  Stock keeping unit: %s\n", a.IdentityData.Sku)
		}
		if a.IdentityData.Sn != "" {
			fmt.Printf("  Serial number: %s\n", a.IdentityData.Sn)
		}
	}
	if detailLevel >= 1 {
		fmt.Printf("CreatedTs: %s\n", a.CreatedTs)
		fmt.Printf("UpdatedTs: %s\n", a.UpdatedTs)
		fmt.Printf("Decommissioning: %t\n", a.Decommissioning)
	}
	if detailLevel >= 2 {
		for i, v := range a.AuthSets {
			fmt.Printf("AuthSet[%d]:\n", i)
			fmt.Printf("  ID: %s\n", v.ID)
			fmt.Printf("  PubKey:\n%s", v.PubKey)
			fmt.Println("  IdentityData:")
			if v.IdentityData.Mac != "" {
				fmt.Printf("    MAC address: %s\n", v.IdentityData.Mac)
			}
			if v.IdentityData.Sku != "" {
				fmt.Printf("    Stock keeping unit: %s\n", v.IdentityData.Sku)
			}
			if v.IdentityData.Sn != "" {
				fmt.Printf("    Serial number: %s\n", v.IdentityData.Sn)
			}
			fmt.Printf("  Status: %s\n", v.Status)
			fmt.Printf("  Ts: %s\n", v.Ts)
		}
	}

	fmt.Println("--------------------------------------------------------------------------------")
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
