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
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/mendersoftware/mender-cli/client/inventory"
)

const cliArgsEnv = "MENDER_CLI_TEST_ARGS"

func TestMain(m *testing.M) {
	if args, ok := os.LookupEnv(cliArgsEnv); ok {
		var argv []string
		if err := json.Unmarshal([]byte(args), &argv); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		os.Args = append([]string{"mender-cli"}, argv...)
		Execute()
		os.Exit(0)
	}
	os.Exit(m.Run())
}

type cliResult struct {
	stdout, stderr string
	exitCode       int
}

func runCLI(t *testing.T, args ...string) cliResult {
	argv, err := json.Marshal(args)
	if err != nil {
		t.Fatal(err)
	}
	home := t.TempDir()
	c := exec.Command(os.Args[0])
	c.Dir = home
	c.Env = append(os.Environ(),
		cliArgsEnv+"="+string(argv),
		"HOME="+home,
		"XDG_CACHE_HOME="+filepath.Join(home, ".cache"),
	)
	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr
	res := cliResult{}
	err = c.Run()
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		res.exitCode = exitErr.ExitCode()
	} else if err != nil {
		t.Fatal(err)
	}
	res.stdout, res.stderr = stdout.String(), stderr.String()
	return res
}

type fakeMender struct {
	*httptest.Server
	devices      []inventory.Device
	total        int
	searchStatus int

	mu       sync.Mutex
	searches []inventory.SearchParams
}

func newFakeMender(t *testing.T, n int) *fakeMender {
	f := &fakeMender{total: n}
	for i := 0; i < n; i++ {
		f.devices = append(f.devices, inventory.Device{
			ID: fmt.Sprintf("dev%d", i),
			Attributes: []inventory.Attribute{
				{Name: "hostname", Scope: "inventory", Value: fmt.Sprintf("pos-%d", i)},
				{Name: "device_type", Scope: "inventory", Value: "precision2"},
				{Name: "mac", Scope: "identity", Value: fmt.Sprintf("00:00:00:00:00:0%d", i)},
			},
		})
	}
	f.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		switch {
		case r.Method == http.MethodPost &&
			r.URL.Path == "/api/management/v2/inventory/filters/search":
			f.search(t, w, r)
		case r.Method == http.MethodGet &&
			r.URL.Path == "/api/management/v2/inventory/filters/attributes":
			_ = json.NewEncoder(w).Encode([]inventory.FilterAttribute{
				{Name: "hostname", Scope: "inventory", Count: n},
				{Name: "mac", Scope: "identity", Count: n},
			})
		case r.Method == http.MethodGet &&
			strings.HasPrefix(r.URL.Path, "/api/management/v1/inventory/devices/"):
			id := strings.TrimPrefix(r.URL.Path, "/api/management/v1/inventory/devices/")
			for _, dev := range f.devices {
				if dev.ID == id {
					_ = json.NewEncoder(w).Encode(dev)
					return
				}
			}
			w.WriteHeader(http.StatusNotFound)
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(f.Close)
	return f
}

func (f *fakeMender) search(t *testing.T, w http.ResponseWriter, r *http.Request) {
	var params inventory.SearchParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		t.Errorf("invalid search request: %s", err)
	}
	f.mu.Lock()
	f.searches = append(f.searches, params)
	f.mu.Unlock()

	if f.searchStatus != 0 {
		w.WriteHeader(f.searchStatus)
		_, _ = w.Write([]byte(`{"error": "invalid filter", "request_id": "req1"}`))
		return
	}
	devices := []inventory.Device{}
	for i := (params.Page - 1) * params.PerPage; i < params.Page*params.PerPage &&
		i < len(f.devices); i++ {
		devices = append(devices, f.devices[i])
	}
	w.Header().Set("X-Total-Count", strconv.Itoa(f.total))
	_ = json.NewEncoder(w).Encode(devices)
}

func (f *fakeMender) requests() []inventory.SearchParams {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.searches
}

func inventoryArgs(server string, args ...string) []string {
	return append([]string{"--server", server, "--token-value", "token",
		"devices", "inventory"}, args...)
}

func TestDevicesInventoryList(t *testing.T) {
	f := newFakeMender(t, 3)

	res := runCLI(t, inventoryArgs(f.URL, "list")...)

	if res.exitCode != 0 {
		t.Fatalf("exit code %d: %s", res.exitCode, res.stderr)
	}
	for _, want := range []string{"ID: dev0\n", "inventory/hostname: pos-2\n",
		"identity/mac: 00:00:00:00:00:01\n"} {
		if !strings.Contains(res.stdout, want) {
			t.Errorf("output does not contain %q:\n%s", want, res.stdout)
		}
	}
	reqs := f.requests()
	if len(reqs) != 1 || reqs[0].Page != 1 || reqs[0].PerPage != 20 {
		t.Fatalf("unexpected requests: %+v", reqs)
	}
	wantSort := []inventory.SortCriteria{
		{Scope: "system", Attribute: "created_ts", Order: "asc"},
		{Scope: "identity", Attribute: "id", Order: "asc"},
	}
	if !reflect.DeepEqual(reqs[0].Sort, wantSort) || reqs[0].Attributes != nil {
		t.Errorf("unexpected request: %+v", reqs[0])
	}
}

func TestDevicesInventoryListFilters(t *testing.T) {
	f := newFakeMender(t, 2)

	res := runCLI(t, inventoryArgs(f.URL, "list",
		"--filter", "device_type=in:precision2,spectrum",
		"--filter", "identity/status=accepted",
		"--filter", "hostname=~pos-",
		"--filter", "!tags/broken?",
		"--columns", "hostname,mac",
		"--sort", "hostname:desc",
		"--output", "csv")...)

	if res.exitCode != 0 {
		t.Fatalf("exit code %d: %s", res.exitCode, res.stderr)
	}
	wantOut := "id,hostname,mac\n" +
		"dev0,pos-0,00:00:00:00:00:00\n" +
		"dev1,pos-1,00:00:00:00:00:01\n"
	if res.stdout != wantOut {
		t.Errorf("got:\n%s\nwant:\n%s", res.stdout, wantOut)
	}

	reqs := f.requests()
	if len(reqs) != 1 {
		t.Fatalf("expected one request, got %+v", reqs)
	}
	filters, err := json.Marshal(reqs[0].Filters)
	if err != nil {
		t.Fatal(err)
	}
	wantFilters := `[` +
		`{"scope":"inventory","attribute":"device_type","type":"$in",` +
		`"value":["precision2","spectrum"]},` +
		`{"scope":"identity","attribute":"status","type":"$eq","value":"accepted"},` +
		`{"scope":"inventory","attribute":"hostname","type":"$regex","value":"pos-"},` +
		`{"scope":"tags","attribute":"broken","type":"$exists","value":false}]`
	if string(filters) != wantFilters {
		t.Errorf("got filters:\n%s\nwant:\n%s", filters, wantFilters)
	}
	if reqs[0].Sort[0] != (inventory.SortCriteria{
		Scope: "inventory", Attribute: "hostname", Order: "desc"}) {
		t.Errorf("unexpected sort: %+v", reqs[0].Sort)
	}
	selected := map[inventory.SelectAttribute]bool{}
	for _, a := range reqs[0].Attributes {
		selected[a] = true
	}
	for _, a := range []inventory.SelectAttribute{
		{Scope: "inventory", Attribute: "hostname"},
		{Scope: "identity", Attribute: "mac"},
		{Scope: "system", Attribute: "created_ts"},
	} {
		if !selected[a] {
			t.Errorf("attribute %+v is not selected: %+v", a, reqs[0].Attributes)
		}
	}
}

func TestDevicesInventoryListAll(t *testing.T) {
	f := newFakeMender(t, 5)
	outFile := filepath.Join(t.TempDir(), "devices.csv")

	res := runCLI(t, inventoryArgs(f.URL, "list", "--all", "--per-page", "2",
		"--progress", "--columns", "hostname", "--output", "csv",
		"--output-file", outFile)...)

	if res.exitCode != 0 {
		t.Fatalf("exit code %d: %s", res.exitCode, res.stderr)
	}
	if res.stdout != "" {
		t.Errorf("nothing should be written to stdout, got %q", res.stdout)
	}
	for _, want := range []string{"Fetched 2/5 devices (40%)\n",
		"Fetched 5/5 devices (100%)\n"} {
		if !strings.Contains(res.stderr, want) {
			t.Errorf("stderr does not contain %q:\n%s", want, res.stderr)
		}
	}
	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatal(err)
	}
	want := "id,hostname\ndev0,pos-0\ndev1,pos-1\ndev2,pos-2\ndev3,pos-3\ndev4,pos-4\n"
	if string(data) != want {
		t.Errorf("got:\n%s\nwant:\n%s", data, want)
	}
	if entries, _ := os.ReadDir(filepath.Dir(outFile)); len(entries) != 1 {
		t.Errorf("temporary files left behind: %v", entries)
	}
	if reqs := f.requests(); len(reqs) != 3 {
		t.Errorf("expected 3 page requests, got %d", len(reqs))
	}
}

func TestDevicesInventoryListIncomplete(t *testing.T) {
	f := newFakeMender(t, 3)
	f.total = 4

	res := runCLI(t, inventoryArgs(f.URL, "list", "--all", "--per-page", "2",
		"--columns", "hostname", "--output", "csv")...)

	if res.exitCode != exitIncomplete {
		t.Fatalf("expected exit code %d, got %d: %s", exitIncomplete, res.exitCode, res.stderr)
	}
	if !strings.Contains(res.stderr,
		"WARNING: fetched 3 of 4 devices, the inventory changed while fetching") {
		t.Errorf("unexpected stderr: %s", res.stderr)
	}
	if res.stdout != "id,hostname\ndev0,pos-0\ndev1,pos-1\ndev2,pos-2\n" {
		t.Errorf("the output should still be written, got %q", res.stdout)
	}
}

func TestDevicesInventoryListCount(t *testing.T) {
	f := newFakeMender(t, 3)

	res := runCLI(t, inventoryArgs(f.URL, "list", "--count",
		"--filter", "device_type=precision2")...)

	if res.exitCode != 0 || res.stdout != "3\n" {
		t.Fatalf("exit code %d, stdout %q: %s", res.exitCode, res.stdout, res.stderr)
	}
	if reqs := f.requests(); len(reqs) != 1 || reqs[0].PerPage != 1 {
		t.Errorf("unexpected requests: %+v", reqs)
	}
}

func TestDevicesInventoryListErrors(t *testing.T) {
	testCases := map[string]struct {
		args   []string
		status int
		stderr string
	}{
		"server error": {
			args:   []string{"list"},
			status: http.StatusBadRequest,
			stderr: "FAILURE: POST /api/management/v2/inventory/filters/search " +
				"request failed with status 400: invalid filter (request ID: req1)",
		},
		"invalid filter": {
			args:   []string{"list", "--filter", "hostname"},
			stderr: `FAILURE: invalid filter "hostname"`,
		},
		"invalid output": {
			args:   []string{"list", "--output", "xml"},
			stderr: `FAILURE: invalid output format "xml"`,
		},
		"progress without all": {
			args:   []string{"list", "--progress"},
			stderr: "FAILURE: --progress requires --all",
		},
		"count with all": {
			args:   []string{"list", "--count", "--all"},
			stderr: "FAILURE: if any flags in the group [count all] are set",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			f := newFakeMender(t, 1)
			f.searchStatus = tc.status

			res := runCLI(t, inventoryArgs(f.URL, tc.args...)...)

			if res.exitCode != 1 {
				t.Errorf("expected exit code 1, got %d", res.exitCode)
			}
			if !strings.Contains(res.stderr, tc.stderr) {
				t.Errorf("stderr does not contain %q:\n%s", tc.stderr, res.stderr)
			}
			if tc.status == 0 && len(f.requests()) != 0 {
				t.Errorf("no request should be sent, got %+v", f.requests())
			}
		})
	}
}

func TestDevicesInventoryNoLogin(t *testing.T) {
	f := newFakeMender(t, 1)

	for _, args := range [][]string{{"list"}, {"show", "dev0"}, {"attributes"}} {
		res := runCLI(t, append([]string{"--server", f.URL, "devices", "inventory"},
			args...)...)

		if res.exitCode != 1 || !strings.Contains(res.stderr, "Please Login first") {
			t.Errorf("%v: exit code %d: %s", args, res.exitCode, res.stderr)
		}
	}
}

func TestDevicesInventoryShow(t *testing.T) {
	f := newFakeMender(t, 2)

	res := runCLI(t, inventoryArgs(f.URL, "show", "dev1")...)
	if res.exitCode != 0 {
		t.Fatalf("exit code %d: %s", res.exitCode, res.stderr)
	}
	for _, want := range []string{"ID: dev1\n", "inventory/hostname: pos-1\n"} {
		if !strings.Contains(res.stdout, want) {
			t.Errorf("output does not contain %q:\n%s", want, res.stdout)
		}
	}

	res = runCLI(t, inventoryArgs(f.URL, "show", "dev1", "--raw")...)
	var dev inventory.Device
	if err := json.Unmarshal([]byte(res.stdout), &dev); err != nil || dev.ID != "dev1" {
		t.Errorf("unexpected raw output %q: %v", res.stdout, err)
	}

	res = runCLI(t, inventoryArgs(f.URL, "show", "missing")...)
	if res.exitCode != 1 || !strings.Contains(res.stderr, "status 404") {
		t.Errorf("exit code %d: %s", res.exitCode, res.stderr)
	}
}

func TestDevicesInventoryAttributes(t *testing.T) {
	f := newFakeMender(t, 2)

	res := runCLI(t, inventoryArgs(f.URL, "attributes")...)
	if res.exitCode != 0 {
		t.Fatalf("exit code %d: %s", res.exitCode, res.stderr)
	}
	for _, want := range []string{"SCOPE", "inventory  hostname  2", "identity   mac       2"} {
		if !strings.Contains(res.stdout, want) {
			t.Errorf("output does not contain %q:\n%s", want, res.stdout)
		}
	}

	res = runCLI(t, inventoryArgs(f.URL, "attributes", "--output", "csv")...)
	if res.stdout != "scope,name,count\ninventory,hostname,2\nidentity,mac,2\n" {
		t.Errorf("unexpected output %q", res.stdout)
	}
}
