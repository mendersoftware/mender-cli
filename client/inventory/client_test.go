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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
)

func fakeInventory(t *testing.T, n int, requests *[]SearchParams) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != searchURL {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.Header.Get("Authorization") != "Bearer token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var params SearchParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			t.Errorf("invalid request body: %s", err)
		}
		*requests = append(*requests, params)
		devices := []Device{}
		for i := (params.Page - 1) * params.PerPage; i < params.Page*params.PerPage && i < n; i++ {
			devices = append(devices, Device{
				ID:         fmt.Sprintf("dev%d", i),
				Attributes: []Attribute{{Name: "n", Scope: "inventory", Value: i}},
			})
		}
		w.Header().Set("X-Total-Count", strconv.Itoa(n))
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(devices)
	}))
}

func TestSearch(t *testing.T) {
	t.Parallel()
	var requests []SearchParams
	srv := fakeInventory(t, 5, &requests)
	defer srv.Close()

	params := SearchParams{
		Page:    1,
		PerPage: 2,
		Filters: []FilterPredicate{{"inventory", "device_type", "$eq", "precision2"}},
	}
	res, err := NewClient(srv.URL, false).Search("token", params)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if res.Total != 5 || len(res.Devices) != 2 || len(res.Raw) != 2 {
		t.Errorf("unexpected result: %+v", res)
	}
	if res.Devices[1].Attributes[0].Value != json.Number("1") {
		t.Errorf("numbers should be decoded as json.Number, got %#v",
			res.Devices[1].Attributes[0].Value)
	}
	if len(requests) != 1 || requests[0].Filters[0].Attribute != "device_type" {
		t.Errorf("unexpected requests: %+v", requests)
	}
}

func searchEach(t *testing.T, n, perPage int) ([]int, int, int, int) {
	var requests []SearchParams
	srv := fakeInventory(t, n, &requests)
	defer srv.Close()

	var pages []int
	stats, err := NewClient(srv.URL, false).SearchEach("token",
		SearchParams{Page: 1, PerPage: perPage},
		func(page *SearchResult, _ SearchStats) error {
			if len(page.Devices) != len(page.Raw) {
				t.Errorf("devices and raw devices differ: %+v", page)
			}
			pages = append(pages, len(page.Devices))
			return nil
		})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	return pages, stats.Fetched, stats.Total, len(requests)
}

func TestSearchEach(t *testing.T) {
	t.Parallel()
	pages, fetched, total, requests := searchEach(t, 5, 2)
	if !reflect.DeepEqual(pages, []int{2, 2, 1}) || fetched != 5 || total != 5 || requests != 3 {
		t.Errorf("got pages %v, fetched %d of %d in %d requests", pages, fetched, total, requests)
	}
}

func TestSearchEachExactPages(t *testing.T) {
	t.Parallel()
	pages, fetched, total, requests := searchEach(t, 4, 2)
	if !reflect.DeepEqual(pages, []int{2, 2}) || fetched != 4 || total != 4 || requests != 2 {
		t.Errorf("got pages %v, fetched %d of %d in %d requests", pages, fetched, total, requests)
	}
}

func TestSearchEachNoDevices(t *testing.T) {
	t.Parallel()
	pages, fetched, total, requests := searchEach(t, 0, 2)
	if !reflect.DeepEqual(pages, []int{0}) || fetched != 0 || total != 0 || requests != 1 {
		t.Errorf("got pages %v, fetched %d of %d in %d requests", pages, fetched, total, requests)
	}
}

func TestSearchEachSkipsDuplicates(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params SearchParams
		_ = json.NewDecoder(r.Body).Decode(&params)
		devices := []Device{}
		start := (params.Page-1)*params.PerPage - (params.Page - 1)
		for i := start; i < start+params.PerPage && i < 5; i++ {
			devices = append(devices, Device{ID: fmt.Sprintf("dev%d", i)})
		}
		w.Header().Set("X-Total-Count", "10")
		_ = json.NewEncoder(w).Encode(devices)
	}))
	defer srv.Close()

	var ids []string
	stats, err := NewClient(srv.URL, false).SearchEach("token",
		SearchParams{Page: 1, PerPage: 2},
		func(page *SearchResult, _ SearchStats) error {
			for _, dev := range page.Devices {
				ids = append(ids, dev.ID)
			}
			return nil
		})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	want := []string{"dev0", "dev1", "dev2", "dev3", "dev4"}
	if !reflect.DeepEqual(ids, want) || stats != (SearchStats{Fetched: 5, Duplicates: 4, Total: 10}) {
		t.Errorf("got %v, %+v", ids, stats)
	}
}

func TestSearchEachCallbackError(t *testing.T) {
	t.Parallel()
	var requests []SearchParams
	srv := fakeInventory(t, 5, &requests)
	defer srv.Close()

	_, err := NewClient(srv.URL, false).SearchEach("token",
		SearchParams{Page: 1, PerPage: 2},
		func(page *SearchResult, _ SearchStats) error { return errors.New("disk full") })
	if err == nil || err.Error() != "disk full" || len(requests) != 1 {
		t.Errorf("expected the callback error after one request, got %v", err)
	}
}

func TestSearchError(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error": "invalid filter", "request_id": "abc"}`))
	}))
	defer srv.Close()

	_, err := NewClient(srv.URL, false).Search("token", SearchParams{Page: 1, PerPage: 1})
	if err == nil || !strings.Contains(err.Error(), "status 400: invalid filter (request ID: abc)") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRetryOnRateLimit(t *testing.T) {
	t.Parallel()
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&calls, 1) < 3 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		_, _ = w.Write([]byte(`[{"name": "hostname", "scope": "inventory", "count": 3}]`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, false)
	c.retryDelay = 0
	attrs, _, err := c.FilterableAttributes("token")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(attrs) != 1 || attrs[0].Name != "hostname" || calls != 3 {
		t.Errorf("unexpected result %+v after %d calls", attrs, calls)
	}
}

func TestGetDevice(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != deviceURL+"dev1" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, _ = w.Write([]byte(
			`{"id": "dev1", "attributes": [{"name": "a", "scope": "tags", "value": "b"}]}`))
	}))
	defer srv.Close()

	dev, raw, err := NewClient(srv.URL, false).GetDevice("token", "dev1")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if dev.ID != "dev1" || len(dev.Attributes) != 1 || !json.Valid(raw) {
		t.Errorf("unexpected device %+v", dev)
	}
	if _, _, err := NewClient(srv.URL, false).GetDevice("token", "nope"); err == nil {
		t.Error("expected an error for a missing device")
	}
}
