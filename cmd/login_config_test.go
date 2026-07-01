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
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestPersistServerToConfigCreatesNewFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".mender-clirc")

	written, err := persistServerToConfig(path, "https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !written {
		t.Fatal("expected written=true for a new file")
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("failed to stat created file: %s", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("expected file mode 0600, got %o", perm)
	}

	content := readJSONFile(t, path)
	if len(content) != 1 {
		t.Errorf("expected 1 key, got %d: %v", len(content), content)
	}
	if content["server"] != "https://example.com" {
		t.Errorf("expected server=https://example.com, got %v", content["server"])
	}
}

func TestPersistServerToConfigCreatesParentDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", ".mender-clirc")

	written, err := persistServerToConfig(path, "https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !written {
		t.Fatal("expected written=true for a new file")
	}

	content := readJSONFile(t, path)
	if content["server"] != "https://example.com" {
		t.Errorf("expected server=https://example.com, got %v", content["server"])
	}
}

func TestPersistServerToConfigAddsServerPreservingExistingKeys(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".mender-clirc")

	original := []byte(`{"username": "alice", "password": "secret"}`)
	if err := os.WriteFile(path, original, 0644); err != nil {
		t.Fatalf("failed to write fixture file: %s", err)
	}

	written, err := persistServerToConfig(path, "https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !written {
		t.Fatal("expected written=true when server key was absent")
	}

	content := readJSONFile(t, path)
	if len(content) != 3 {
		t.Errorf("expected 3 keys, got %d: %v", len(content), content)
	}
	if content["server"] != "https://example.com" {
		t.Errorf("expected server=https://example.com, got %v", content["server"])
	}
	if content["username"] != "alice" {
		t.Errorf("expected username=alice, got %v", content["username"])
	}
	if content["password"] != "secret" {
		t.Errorf("expected password=secret, got %v", content["password"])
	}
}

func TestPersistServerToConfigDoesNotOverwriteExistingServerKey(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".mender-clirc")

	original := []byte(`{"server": "https://staging.example.com", "username": "bob"}`)
	if err := os.WriteFile(path, original, 0644); err != nil {
		t.Fatalf("failed to write fixture file: %s", err)
	}

	written, err := persistServerToConfig(path, "https://different.example.com")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if written {
		t.Fatal("expected written=false when server key already exists")
	}

	content := readJSONFile(t, path)
	if content["server"] != "https://staging.example.com" {
		t.Errorf("expected server to remain unchanged, got %v", content["server"])
	}
	if content["username"] != "bob" {
		t.Errorf("expected username=bob, got %v", content["username"])
	}
}

func TestPersistServerToConfigEmptyStringServerKeyCountsAsSet(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".mender-clirc")

	original := []byte(`{"server": ""}`)
	if err := os.WriteFile(path, original, 0644); err != nil {
		t.Fatalf("failed to write fixture file: %s", err)
	}

	written, err := persistServerToConfig(path, "https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if written {
		t.Fatal("expected written=false when server key is present (even if empty)")
	}

	content := readJSONFile(t, path)
	if content["server"] != "" {
		t.Errorf("expected server to remain empty, got %v", content["server"])
	}
}

func TestPersistServerToConfigMalformedJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".mender-clirc")

	original := []byte(`{not valid json`)
	if err := os.WriteFile(path, original, 0644); err != nil {
		t.Fatalf("failed to write fixture file: %s", err)
	}

	_, err := persistServerToConfig(path, "https://example.com")
	if err == nil {
		t.Fatal("expected error for malformed JSON, got nil")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %s", err)
	}
	if string(data) != string(original) {
		t.Errorf("expected file to remain unchanged, got %q", string(data))
	}
}

func TestResolveConfigFilePath(t *testing.T) {
	homeErr := errors.New("no home dir")

	cases := []struct {
		name     string
		used     string
		homeDir  string
		homeErr  error
		wantPath string
	}{
		{
			name:     "viper config file used is returned as-is",
			used:     "/etc/mender-cli/.mender-clirc",
			homeDir:  "/home/someone",
			wantPath: "/etc/mender-cli/.mender-clirc",
		},
		{
			name:     "falls back to home dir when no config file used",
			used:     "",
			homeDir:  "/home/someone",
			wantPath: filepath.Join("/home/someone", ".mender-clirc"),
		},
		{
			name:     "falls back to relative path when home dir errors",
			used:     "",
			homeErr:  homeErr,
			wantPath: ".mender-clirc",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			homeDirFunc := func() (string, error) {
				if tc.homeErr != nil {
					return "", tc.homeErr
				}
				return tc.homeDir, nil
			}

			got := resolveConfigFilePath(tc.used, homeDirFunc)
			if got != tc.wantPath {
				t.Errorf("resolveConfigFilePath(%q, ...) = %q, want %q", tc.used, got, tc.wantPath)
			}
		})
	}
}

func readJSONFile(t *testing.T, path string) map[string]interface{} {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file %s: %s", path, err)
	}

	var content map[string]interface{}
	if err := json.Unmarshal(data, &content); err != nil {
		t.Fatalf("failed to parse JSON file %s: %s", path, err)
	}

	return content
}
