// Copyright 2026 Northern.tech AS
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
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"

	"github.com/mendersoftware/mender-cli/client"
	"github.com/mendersoftware/mender-cli/log"
)

const (
	searchURL     = "/api/management/v2/inventory/filters/search"
	attributesURL = "/api/management/v2/inventory/filters/attributes"
	deviceURL     = "/api/management/v1/inventory/devices/"

	requestTimeout = 60 * time.Second
	maxRetries     = 5
)

const (
	ScopeInventory = "inventory"
	ScopeIdentity  = "identity"
	ScopeSystem    = "system"
	ScopeTags      = "tags"
	ScopeMonitor   = "monitor"
)

var Scopes = []string{ScopeInventory, ScopeIdentity, ScopeSystem, ScopeTags, ScopeMonitor}

type Attribute struct {
	Name        string      `json:"name"`
	Scope       string      `json:"scope"`
	Value       interface{} `json:"value"`
	Description string      `json:"description,omitempty"`
	Timestamp   string      `json:"timestamp,omitempty"`
}

type Device struct {
	ID         string      `json:"id"`
	UpdatedTs  string      `json:"updated_ts,omitempty"`
	Attributes []Attribute `json:"attributes"`
}

type FilterPredicate struct {
	Scope     string      `json:"scope"`
	Attribute string      `json:"attribute"`
	Type      string      `json:"type"`
	Value     interface{} `json:"value"`
}

type SortCriteria struct {
	Scope     string `json:"scope"`
	Attribute string `json:"attribute"`
	Order     string `json:"order"`
}

type SelectAttribute struct {
	Scope     string `json:"scope"`
	Attribute string `json:"attribute"`
}

type SearchParams struct {
	Page       int               `json:"page"`
	PerPage    int               `json:"per_page"`
	Filters    []FilterPredicate `json:"filters,omitempty"`
	Sort       []SortCriteria    `json:"sort,omitempty"`
	Attributes []SelectAttribute `json:"attributes,omitempty"`
}

type FilterAttribute struct {
	Name  string `json:"name"`
	Scope string `json:"scope"`
	Count int    `json:"count"`
}

type SearchResult struct {
	Devices []Device
	Raw     []json.RawMessage
	Total   int
}

type apiError struct {
	Error     string `json:"error"`
	RequestID string `json:"request_id"`
}

type Client struct {
	searchURL     string
	attributesURL string
	deviceURL     string
	client        *http.Client
	retryDelay    time.Duration
}

func NewClient(url string, skipVerify bool) *Client {
	httpClient := client.NewHttpClient(skipVerify)
	httpClient.Timeout = requestTimeout
	return &Client{
		searchURL:     client.JoinURL(url, searchURL),
		attributesURL: client.JoinURL(url, attributesURL),
		deviceURL:     client.JoinURL(url, deviceURL),
		client:        httpClient,
		retryDelay:    time.Second,
	}
}

func (c *Client) Search(token string, params SearchParams) (*SearchResult, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	hdr, data, err := c.do(token, http.MethodPost, c.searchURL, body)
	if err != nil {
		return nil, err
	}
	res := &SearchResult{}
	if err := json.Unmarshal(data, &res.Raw); err != nil {
		return nil, fmt.Errorf("failed to parse search response: %w", err)
	}
	res.Devices = make([]Device, len(res.Raw))
	for i, raw := range res.Raw {
		if err := decodeDevice(raw, &res.Devices[i]); err != nil {
			return nil, err
		}
	}
	res.Total = len(res.Devices)
	if total := hdr.Get("X-Total-Count"); total != "" {
		if res.Total, err = strconv.Atoi(total); err != nil {
			return nil, fmt.Errorf("invalid X-Total-Count header %q: %w", total, err)
		}
	}
	return res, nil
}

type SearchStats struct {
	Fetched    int
	Duplicates int
	Total      int
}

func (c *Client) SearchEach(
	token string,
	params SearchParams,
	fn func(page *SearchResult, stats SearchStats) error,
) (SearchStats, error) {
	var stats SearchStats
	seen := map[string]bool{}
	for {
		res, err := c.Search(token, params)
		if err != nil {
			return stats, err
		}
		stats.Total = res.Total
		page := &SearchResult{Total: res.Total}
		for i, dev := range res.Devices {
			if seen[dev.ID] {
				stats.Duplicates++
				continue
			}
			seen[dev.ID] = true
			page.Devices = append(page.Devices, dev)
			page.Raw = append(page.Raw, res.Raw[i])
		}
		stats.Fetched += len(page.Devices)
		if err := fn(page, stats); err != nil {
			return stats, err
		}
		if len(res.Devices) < params.PerPage ||
			params.Page*params.PerPage >= res.Total {
			return stats, nil
		}
		params.Page++
	}
}

func (c *Client) GetDevice(token, id string) (*Device, json.RawMessage, error) {
	_, data, err := c.do(token, http.MethodGet, c.deviceURL+url.PathEscape(id), nil)
	if err != nil {
		return nil, nil, err
	}
	var dev Device
	if err := decodeDevice(data, &dev); err != nil {
		return nil, nil, err
	}
	return &dev, data, nil
}

func (c *Client) FilterableAttributes(token string) ([]FilterAttribute, []byte, error) {
	_, data, err := c.do(token, http.MethodGet, c.attributesURL, nil)
	if err != nil {
		return nil, nil, err
	}
	var attrs []FilterAttribute
	if err := json.Unmarshal(data, &attrs); err != nil {
		return nil, nil, fmt.Errorf("failed to parse attributes response: %w", err)
	}
	return attrs, data, nil
}

func decodeDevice(data []byte, dev *Device) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := dec.Decode(dev); err != nil {
		return fmt.Errorf("failed to parse device inventory: %w", err)
	}
	return nil
}

func (c *Client) do(token, method, reqURL string, body []byte) (http.Header, []byte, error) {
	for attempt := 0; ; attempt++ {
		req, err := http.NewRequest(method, reqURL, bytes.NewReader(body))
		if err != nil {
			return nil, nil, fmt.Errorf("failed to prepare request: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		reqDump, err := httputil.DumpRequest(req, false)
		if err != nil {
			return nil, nil, err
		}
		log.Verbf("sending request: \n%s", string(reqDump))
		if body != nil {
			log.Verbf("request body: %s", string(body))
		}

		rsp, err := c.client.Do(req)
		if err != nil {
			return nil, nil, fmt.Errorf("%s %s request failed: %w", method, reqURL, err)
		}
		data, err := io.ReadAll(rsp.Body)
		rsp.Body.Close()
		if err != nil {
			return nil, nil, fmt.Errorf("error reading response body: %w", err)
		}

		if rsp.StatusCode == http.StatusTooManyRequests && attempt < maxRetries {
			delay := c.retryDelay << attempt
			if s, err := strconv.Atoi(rsp.Header.Get("Retry-After")); err == nil {
				delay = time.Duration(s) * time.Second
			}
			log.Verbf("rate limited, retrying in %s", delay)
			time.Sleep(delay)
			continue
		}
		if rsp.StatusCode != http.StatusOK {
			msg := fmt.Sprintf("%s %s request failed with status %d",
				method, req.URL.RequestURI(), rsp.StatusCode)
			var apiErr apiError
			if json.Unmarshal(data, &apiErr) == nil && apiErr.Error != "" {
				msg += ": " + apiErr.Error
				if apiErr.RequestID != "" {
					msg += " (request ID: " + apiErr.RequestID + ")"
				}
			}
			return nil, nil, fmt.Errorf("%s", msg)
		}
		return rsp.Header, data, nil
	}
}
