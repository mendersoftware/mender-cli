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
package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/mendersoftware/mender-cli/client/inventory"
)

func fakeInventory(t *testing.T, n, total, failPage int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params inventory.SearchParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			t.Errorf("invalid request body: %s", err)
		}
		if params.Page == failPage {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		devices := []inventory.Device{}
		for i := (params.Page - 1) * params.PerPage; i < params.Page*params.PerPage && i < n; i++ {
			devices = append(devices, inventory.Device{
				ID: fmt.Sprintf("dev%d", i),
				Attributes: []inventory.Attribute{
					{Name: "hostname", Scope: "inventory", Value: fmt.Sprintf("host%d", i)},
				},
			})
		}
		w.Header().Set("X-Total-Count", strconv.Itoa(total))
		_ = json.NewEncoder(w).Encode(devices)
	}))
}

func newTestListCmd(server string) (*DevicesInventoryListCmd, *bytes.Buffer, *bytes.Buffer) {
	cols, _ := inventory.ParseColumns([]string{"hostname"})
	var stdout, stderr bytes.Buffer
	return &DevicesInventoryListCmd{
		inventoryCmdCommon: &inventoryCmdCommon{
			server:  server,
			token:   "token",
			output:  inventory.FormatCSV,
			columns: cols,
		},
		params: inventory.SearchParams{Page: 1, PerPage: 2},
		all:    true,
		stdout: &stdout,
		stderr: &stderr,
	}, &stdout, &stderr
}

func dirEntries(t *testing.T, dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	return names
}

func TestInventoryListAllWithProgress(t *testing.T) {
	srv := fakeInventory(t, 5, 5, 0)
	defer srv.Close()
	c, stdout, stderr := newTestListCmd(srv.URL)
	c.progress = true
	c.outputFile = filepath.Join(t.TempDir(), "out.csv")

	if err := c.Run(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	data, err := os.ReadFile(c.outputFile)
	if err != nil {
		t.Fatal(err)
	}
	want := "id,hostname\ndev0,host0\ndev1,host1\ndev2,host2\ndev3,host3\ndev4,host4\n"
	if string(data) != want {
		t.Errorf("got:\n%s\nwant:\n%s", data, want)
	}
	if stdout.Len() != 0 {
		t.Errorf("nothing should be written to stdout, got %q", stdout.String())
	}
	wantProgress := "Fetched 2/5 devices (40%)\n" +
		"Fetched 4/5 devices (80%)\n" +
		"Fetched 5/5 devices (100%)\n"
	if stderr.String() != wantProgress {
		t.Errorf("got progress:\n%s\nwant:\n%s", stderr.String(), wantProgress)
	}
	if names := dirEntries(t, filepath.Dir(c.outputFile)); len(names) != 1 {
		t.Errorf("temporary files left behind: %v", names)
	}
}

func TestInventoryListOutputFileNotReplacedOnError(t *testing.T) {
	srv := fakeInventory(t, 5, 5, 2)
	defer srv.Close()
	c, _, _ := newTestListCmd(srv.URL)
	c.outputFile = filepath.Join(t.TempDir(), "out.csv")
	if err := os.WriteFile(c.outputFile, []byte("previous\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := c.Run()
	if err == nil || !strings.Contains(err.Error(), "status 500") {
		t.Fatalf("expected the server error, got %v", err)
	}
	data, _ := os.ReadFile(c.outputFile)
	if string(data) != "previous\n" {
		t.Errorf("the output file should be untouched, got %q", data)
	}
	if names := dirEntries(t, filepath.Dir(c.outputFile)); len(names) != 1 {
		t.Errorf("temporary files left behind: %v", names)
	}
}

func TestInventoryListIncomplete(t *testing.T) {
	srv := fakeInventory(t, 5, 7, 0)
	defer srv.Close()
	c, _, _ := newTestListCmd(srv.URL)
	c.outputFile = filepath.Join(t.TempDir(), "out.csv")

	err := c.Run()
	var exitErr *exitError
	if !errors.As(err, &exitErr) || exitErr.code != exitIncomplete {
		t.Fatalf("expected exit code %d, got %v", exitIncomplete, err)
	}
	if exitErr.Error() != "fetched 5 of 7 devices, the inventory changed while fetching" {
		t.Errorf("unexpected message: %s", exitErr)
	}
	data, _ := os.ReadFile(c.outputFile)
	if strings.Count(string(data), "\n") != 6 {
		t.Errorf("the output should still be written, got %q", data)
	}
}

func TestInventoryListDuplicates(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Total-Count", "4")
		_ = json.NewEncoder(w).Encode([]inventory.Device{{ID: "dev0"}, {ID: "dev1"}})
	}))
	defer srv.Close()
	c, stdout, _ := newTestListCmd(srv.URL)

	err := c.Run()
	var exitErr *exitError
	if !errors.As(err, &exitErr) || exitErr.code != exitIncomplete {
		t.Fatalf("expected exit code %d, got %v", exitIncomplete, err)
	}
	want := "fetched 2 of 4 devices, 2 were returned more than once; " +
		"the server returned inconsistent pages or the inventory changed while fetching"
	if exitErr.Error() != want {
		t.Errorf("unexpected message: %s", exitErr)
	}
	if stdout.String() != "id,hostname\ndev0,\ndev1,\n" {
		t.Errorf("unexpected output %q", stdout.String())
	}
}

func TestInventoryListAllFromLaterPage(t *testing.T) {
	srv := fakeInventory(t, 5, 5, 0)
	defer srv.Close()
	c, stdout, stderr := newTestListCmd(srv.URL)
	c.params.Page = 2
	c.progress = true

	if err := c.Run(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if stdout.String() != "id,hostname\ndev2,host2\ndev3,host3\ndev4,host4\n" {
		t.Errorf("unexpected output %q", stdout.String())
	}
	if !strings.HasSuffix(stderr.String(), "Fetched 3/3 devices (100%)\n") {
		t.Errorf("unexpected progress %q", stderr.String())
	}
}

func TestInventoryListSinglePage(t *testing.T) {
	srv := fakeInventory(t, 5, 5, 0)
	defer srv.Close()
	c, stdout, _ := newTestListCmd(srv.URL)
	c.all = false

	if err := c.Run(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if stdout.String() != "id,hostname\ndev0,host0\ndev1,host1\n" {
		t.Errorf("unexpected output %q", stdout.String())
	}
}

func TestInventoryListCount(t *testing.T) {
	srv := fakeInventory(t, 5, 5, 0)
	defer srv.Close()
	c, stdout, _ := newTestListCmd(srv.URL)
	c.all = false
	c.count = true

	if err := c.Run(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if stdout.String() != "5\n" {
		t.Errorf("unexpected output %q", stdout.String())
	}
}

func TestValidateListOptions(t *testing.T) {
	testCases := []struct {
		name       string
		all        bool
		progress   bool
		outputFile string
		terminal   bool
		err        string
	}{
		{name: "no progress on terminal", all: true, terminal: true},
		{name: "progress to file", all: true, progress: true, outputFile: "x", terminal: true},
		{name: "progress to pipe", all: true, progress: true},
		{name: "progress on terminal", all: true, progress: true, terminal: true,
			err: "--progress needs --output-file or stdout redirected"},
		{name: "progress without all", progress: true, outputFile: "x",
			err: "--progress requires --all"},
	}
	for _, tc := range testCases {
		err := validateListOptions(tc.all, tc.progress, tc.outputFile, tc.terminal)
		if tc.err == "" && err != nil {
			t.Errorf("%s: unexpected error: %s", tc.name, err)
		}
		if tc.err != "" && (err == nil || !strings.Contains(err.Error(), tc.err)) {
			t.Errorf("%s: expected error %q, got %v", tc.name, tc.err, err)
		}
	}
}

func TestProgressLine(t *testing.T) {
	for _, tc := range []struct {
		fetched, total int
		want           string
	}{
		{0, 0, "Fetched 0/0 devices (100%)"},
		{1, 3, "Fetched 1/3 devices (33%)"},
		{500, 3285, "Fetched 500/3285 devices (15%)"},
		{4, 3, "Fetched 4/3 devices (100%)"},
	} {
		if got := progressLine(tc.fetched, tc.total); got != tc.want {
			t.Errorf("got %q, want %q", got, tc.want)
		}
	}
}
