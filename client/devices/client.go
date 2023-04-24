// Copyright 2023 Northern.tech AS
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
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mendersoftware/mender-cli/client"
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
	return &Client{
		url:            url,
		devicesListURL: client.JoinURL(url, devicesListURL),
		client:         client.NewHttpClient(skipVerify),
	}
}

func (c *Client) ListDevices(tokenPath string, detailLevel int, raw bool) error {
	if detailLevel > 3 || detailLevel < 0 {
		return fmt.Errorf("FAILURE: invalid devices detail")
	}

	body, err := client.DoGetRequest(tokenPath, c.devicesListURL, c.client)
	if err != nil {
		return err
	}

	if raw {
		fmt.Println(string(body))
	} else {

		var list devicesList
		err = json.Unmarshal(body, &list.devices)
		if err != nil {
			return err
		}
		for _, v := range list.devices {
			listDevice(v, detailLevel)
		}
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
