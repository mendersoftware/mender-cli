// Copyright 2025 Northern.tech AS
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
package inventory

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mendersoftware/mender-cli/client"
	"github.com/mendersoftware/mender-cli/log"
)

type devicesList struct {
	Devices []deviceData
}

type deviceData struct {
	ID        string `json:"id"`
	UpdatedTs string `json:"updated_ts"`
}

type FilterPredicate struct {
	Attribute string `json:"attribute"`
	Scope     string `json:"scope"`
	Type      string `json:"type"`
	Value     string `json:"value"`
}

type InventoryAttribute struct {
	Name  string `json:"name"`
	Scope string `json:"scope"`
}

type SearchDeviceInventoriesRequest struct {
	Filters []FilterPredicate `json:"filters"`
}

const (
	inventorySearchURL     = "/api/management/v2/inventory/filters/search"
	inventoryAttributesURL = "/api/management/v2/inventory/filters/attributes"
)

type Client struct {
	url     string
	token   string
	BaseURL string
	client  *http.Client
}

func NewClient(url string, token string, skipVerify bool) *Client {
	return &Client{
		url:     url,
		token:   token,
		BaseURL: url,
		client:  client.NewHttpClient(skipVerify),
	}
}

func (c *Client) SearchDevices(filters []FilterPredicate) (*devicesList, error) {
	if len(filters) == 0 {
		return nil, fmt.Errorf("FAILURE: specify at least one filter")
	}

	req := SearchDeviceInventoriesRequest{
		Filters: filters,
	}

	log.Verbf("Search Devices request: %s\n", req)

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(req)
	if err != nil {
		log.Err(err.Error())
	}
	url := client.JoinURL(c.BaseURL, inventorySearchURL)
	body, err := client.DoPostRequest(c.token, url, c.client, &buf)
	if err != nil {
		return nil, err
	}

	log.Verbf("Search Devices response: %s\n", body)
	var list devicesList
	err = json.Unmarshal(body, &list.Devices)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func (c *Client) ListAttributes() (*[]InventoryAttribute, error) {
	url := client.JoinURL(c.BaseURL, inventoryAttributesURL)
	body, err := client.DoGetRequest(c.token, url, c.client)
	if err != nil {
		return nil, err
	}

	log.Verbf("Inventory attributes response: %s\n", body)
	var list []InventoryAttribute
	err = json.Unmarshal(body, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}
