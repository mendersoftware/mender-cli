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
	"encoding/base64"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"
)

func makeJWT(t *testing.T, payload map[string]any) string {
	t.Helper()
	header := map[string]any{"alg": "HS256", "typ": "JWT"}
	enc := func(v any) string {
		b, err := json.Marshal(v)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		return base64.RawURLEncoding.EncodeToString(b)
	}
	return enc(header) + "." + enc(payload) + ".sig"
}

func TestDecodeJWT_RoundTrip(t *testing.T) {
	tok := makeJWT(t, map[string]any{"sub": "alice", "exp": float64(123)})
	h, p, err := decodeJWT(tok)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if h["alg"] != "HS256" {
		t.Errorf("header alg = %v", h["alg"])
	}
	if p["sub"] != "alice" {
		t.Errorf("payload sub = %v", p["sub"])
	}
}

func TestDecodeJWT_Errors(t *testing.T) {
	cases := map[string]string{
		"single segment": "abc",
		"two segments":   "abc.def",
		"bad base64":     "!!!.def.sig",
		"bad json":       base64.RawURLEncoding.EncodeToString([]byte("not-json")) + ".def.sig",
	}
	for name, tok := range cases {
		t.Run(name, func(t *testing.T) {
			if _, _, err := decodeJWT(tok); err == nil {
				t.Errorf("expected error for %q", tok)
			}
		})
	}
}

func TestJWTExpiry(t *testing.T) {
	exp := time.Now().Add(7 * 24 * time.Hour).Unix()
	_, p, err := decodeJWT(makeJWT(t, map[string]any{"exp": float64(exp)}))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	got, ok := jwtExpiry(p)
	if !ok {
		t.Fatal("expected exp claim")
	}
	if got.Unix() != exp {
		t.Errorf("exp = %d, want %d", got.Unix(), exp)
	}

	if _, ok := jwtExpiry(map[string]any{}); ok {
		t.Error("expected no exp")
	}
	if _, ok := jwtExpiry(map[string]any{"exp": "not-a-number"}); ok {
		t.Error("expected exp parse failure to return ok=false")
	}
}

func TestHumanDuration(t *testing.T) {
	cases := []struct {
		d    time.Duration
		want string
	}{
		{30 * time.Second, "less than a minute"},
		{time.Minute, "1 minute"},
		{2 * time.Minute, "2 minutes"},
		{time.Hour, "1 hour"},
		{3 * time.Hour, "3 hours"},
		{27 * 24 * time.Hour, "27 days"},
	}
	for _, tc := range cases {
		if got := humanDuration(tc.d); got != tc.want {
			t.Errorf("humanDuration(%s) = %q, want %q", tc.d, got, tc.want)
		}
	}
}

// stripANSI is unused but kept available; the real value is the explicit
// assertions below.

func TestSetTokenCmd_Run_StdinPiped(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/authtoken"
	c := &SetTokenCmd{
		server:     "http://invalid.example",
		tokenPath:  path,
		stdin:      strings.NewReader("  abc.def.sig\n"),
		stdinIsTTY: false,
		prompt:     func() ([]byte, error) { t.Fatal("prompt should not be called"); return nil, nil },
		verifier:   func(string, bool, string) error { return nil },
	}
	if err := c.Run(); err != nil {
		t.Fatalf("Run: %v", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "abc.def.sig" {
		t.Errorf("file contents = %q, want %q", string(got), "abc.def.sig")
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if mode := info.Mode().Perm(); mode != 0600 {
		t.Errorf("mode = %o, want 0600", mode)
	}
}

func TestSetTokenCmd_Run_RejectsEmpty(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/authtoken"
	c := &SetTokenCmd{
		server:     "http://invalid.example",
		tokenPath:  path,
		stdin:      strings.NewReader("   \n"),
		stdinIsTTY: false,
	}
	if err := c.Run(); err == nil {
		t.Fatal("expected error for empty token")
	}
	if _, err := os.Stat(path); err == nil {
		t.Error("file should not be created for empty input")
	}
}
