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
	"strings"
	"testing"
)

func testDevices() []Device {
	return []Device{{
		ID:        "dev1",
		UpdatedTs: "2026-10-01T10:00:00Z",
		Attributes: []Attribute{
			{Name: "mac", Scope: "identity", Value: "aa:aa"},
			{Name: "mac", Scope: "inventory", Value: "bb:bb", Description: "MAC address"},
			{Name: "hostname", Scope: "inventory", Value: "pos-1"},
			{Name: "ipv4", Scope: "inventory", Value: []interface{}{"10.0.0.1", "10.0.0.2"}},
			{Name: "cpus", Scope: "inventory", Value: json.Number("4")},
		},
	}, {
		ID: "dev2",
		Attributes: []Attribute{
			{Name: "mac", Scope: "identity", Value: "cc,cc"},
		},
	}}
}

func TestLookup(t *testing.T) {
	t.Parallel()
	dev := testDevices()[0]
	attr, ok := dev.Lookup(Column{Name: "mac"})
	if !ok || attr.Value != "bb:bb" {
		t.Errorf("unscoped lookup should prefer the inventory scope, got %+v", attr)
	}
	attr, ok = dev.Lookup(Column{Scope: "identity", Name: "mac"})
	if !ok || attr.Value != "aa:aa" {
		t.Errorf("scoped lookup returned %+v", attr)
	}
	if _, ok = dev.Lookup(Column{Name: "missing"}); ok {
		t.Error("lookup of a missing attribute should fail")
	}
}

func TestWriteCSV(t *testing.T) {
	t.Parallel()
	cols, _ := ParseColumns([]string{"hostname", "identity/mac", "ipv4", "cpus"})
	var buf bytes.Buffer
	if err := Write(&buf, FormatCSV, testDevices(), cols, 0); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	want := "id,hostname,identity/mac,ipv4,cpus\n" +
		"dev1,pos-1,aa:aa,10.0.0.1;10.0.0.2,4\n" +
		"dev2,,\"cc,cc\",,\n"
	if buf.String() != want {
		t.Errorf("got:\n%s\nwant:\n%s", buf.String(), want)
	}
}

func TestWriteTSVAllColumns(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	if err := Write(&buf, FormatTSV, testDevices(), nil, 0); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	lines := strings.Split(buf.String(), "\n")
	want := "id\tinventory/cpus\tinventory/hostname\tinventory/ipv4\tinventory/mac\tidentity/mac"
	if lines[0] != want {
		t.Errorf("got header %q, want %q", lines[0], want)
	}
}

func TestWriteJSON(t *testing.T) {
	t.Parallel()
	cols, _ := ParseColumns([]string{"hostname", "ipv4"})
	var buf bytes.Buffer
	if err := Write(&buf, FormatJSON, testDevices(), cols, 0); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	var got []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON output: %s", err)
	}
	if len(got) != 2 || got[0]["id"] != "dev1" || got[0]["hostname"] != "pos-1" {
		t.Errorf("unexpected output: %s", buf.String())
	}
	if ipv4, ok := got[0]["ipv4"].([]interface{}); !ok || len(ipv4) != 2 {
		t.Errorf("array values should be kept as arrays: %s", buf.String())
	}
	if v, ok := got[1]["hostname"]; !ok || v != nil {
		t.Errorf("missing attributes should be null: %s", buf.String())
	}
}

func TestWriteText(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	if err := Write(&buf, FormatText, testDevices()[:1], nil, 1); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	for _, want := range []string{
		"ID: dev1\n",
		"UpdatedTs: 2026-10-01T10:00:00Z\n",
		"inventory/mac: bb:bb\n  Description: MAC address\n",
		"identity/mac: aa:aa\n",
		"inventory/ipv4: 10.0.0.1;10.0.0.2\n",
	} {
		if !strings.Contains(buf.String(), want) {
			t.Errorf("output does not contain %q:\n%s", want, buf.String())
		}
	}
}

func TestWriteRaw(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	if err := WriteRaw(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if strings.TrimSpace(buf.String()) != "[]" {
		t.Errorf("expected an empty array, got %q", buf.String())
	}
}

func TestWriteFilterableAttributes(t *testing.T) {
	t.Parallel()
	attrs := []FilterAttribute{{Name: "device_type", Scope: "inventory", Count: 12}}
	var buf bytes.Buffer
	if err := WriteFilterableAttributes(&buf, FormatCSV, attrs); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if buf.String() != "scope,name,count\ninventory,device_type,12\n" {
		t.Errorf("unexpected output: %q", buf.String())
	}
	if err := WriteFilterableAttributes(&buf, "xml", attrs); err == nil {
		t.Error("expected an error for an invalid format")
	}
}

func TestWriterStreamsWithColumns(t *testing.T) {
	t.Parallel()
	cols, _ := ParseColumns([]string{"hostname"})
	devices := testDevices()
	for _, format := range []string{FormatText, FormatJSON, FormatCSV, FormatTSV} {
		var buf bytes.Buffer
		w, err := NewWriter(&buf, format, cols, 0, false)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if err := w.WritePage(devices[:1], nil); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if !strings.Contains(buf.String(), "dev1") {
			t.Errorf("%s: first page not written before Close: %q", format, buf.String())
		}
		if err := w.WritePage(devices[1:], nil); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if err := w.Close(); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		var single bytes.Buffer
		if err := Write(&single, format, devices, cols, 0); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if buf.String() != single.String() {
			t.Errorf("%s: paged output differs:\n%s\nfrom:\n%s", format, buf.String(), single.String())
		}
	}
}

func TestWriterBuffersWithoutColumns(t *testing.T) {
	t.Parallel()
	devices := testDevices()
	var buf bytes.Buffer
	w, _ := NewWriter(&buf, FormatCSV, nil, 0, false)
	_ = w.WritePage(devices[:1], nil)
	_ = w.WritePage(devices[1:], nil)
	if buf.Len() != 0 {
		t.Errorf("CSV without columns should be written on Close, got %q", buf.String())
	}
	if err := w.Close(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 || !strings.HasSuffix(lines[0], ",identity/mac") {
		t.Errorf("unexpected output:\n%s", buf.String())
	}
}

func TestWriterEmpty(t *testing.T) {
	t.Parallel()
	cols, _ := ParseColumns([]string{"hostname"})
	testCases := []struct {
		format string
		cols   []Column
		raw    bool
		want   string
	}{
		{FormatText, nil, false, ""},
		{FormatJSON, cols, false, "[]\n"},
		{FormatJSON, nil, false, "[]\n"},
		{FormatJSON, nil, true, "[]\n"},
		{FormatCSV, cols, false, "id,hostname\n"},
		{FormatCSV, nil, false, "id\n"},
	}
	for _, tc := range testCases {
		var buf bytes.Buffer
		w, _ := NewWriter(&buf, tc.format, tc.cols, 0, tc.raw)
		_ = w.WritePage(nil, nil)
		if err := w.Close(); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if buf.String() != tc.want {
			t.Errorf("%s (raw %t, cols %d): got %q, want %q",
				tc.format, tc.raw, len(tc.cols), buf.String(), tc.want)
		}
	}
}

func TestWriteRawStreams(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	w, _ := NewWriter(&buf, FormatText, nil, 0, true)
	_ = w.WritePage(make([]Device, 1), []json.RawMessage{json.RawMessage(`{"id":"a"}`)})
	_ = w.WritePage(make([]Device, 1), []json.RawMessage{json.RawMessage(`{"id":"b"}`)})
	_ = w.Close()
	var got []map[string]string
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON %q: %s", buf.String(), err)
	}
	if len(got) != 2 || got[0]["id"] != "a" || got[1]["id"] != "b" {
		t.Errorf("unexpected output %q", buf.String())
	}
}
