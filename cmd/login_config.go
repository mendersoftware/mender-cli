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
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/mendersoftware/mender-cli/log"
)

// persistServerToConfig ensures the JSON config file at path has a "server"
// key set to server. It creates the file (and any missing parent directories)
// if it does not exist. If the file already contains a "server" key, it is
// left untouched and false is returned. Otherwise the key is written and true
// is returned.
func persistServerToConfig(path string, server string) (bool, error) {
	var content map[string]interface{}

	data, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return false, errors.Wrapf(err, "failed to read config file %s", path)
		}
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			return false, errors.Wrapf(err, "failed to create directory for %s", path)
		}
		content = map[string]interface{}{}
	} else {
		if err := json.Unmarshal(data, &content); err != nil {
			return false, errors.Wrapf(err, "failed to parse config file %s as JSON", path)
		}
		if _, exists := content[argRootServer]; exists {
			return false, nil
		}
	}

	content[argRootServer] = server

	out, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		return false, errors.Wrap(err, "failed to encode config file content")
	}
	out = append(out, '\n')

	if err := os.WriteFile(path, out, 0600); err != nil {
		return false, errors.Wrapf(err, "failed to write config file %s", path)
	}

	return true, nil
}

// resolveConfigFilePath returns the config file path to persist the server
// to: the file viper actually loaded, if any, or otherwise
// $HOME/.mender-clirc.
func resolveConfigFilePath(viperConfigFileUsed string, userHomeDir func() (string, error)) string {
	if viperConfigFileUsed != "" {
		return viperConfigFileUsed
	}

	home, err := userHomeDir()
	if err != nil {
		return ".mender-clirc"
	}

	return filepath.Join(home, ".mender-clirc")
}

// maybePersistServer persists c.server to the config file if it was
// explicitly passed via --server and differs from the default server. Any
// error encountered while persisting is logged but not returned, since the
// login itself has already succeeded by the time this is called.
func (c *LoginCmd) maybePersistServer() {
	if !c.serverChanged || c.server == defaultServerURL {
		return
	}

	path := resolveConfigFilePath(viper.ConfigFileUsed(), os.UserHomeDir)
	written, err := persistServerToConfig(path, c.server)
	if err != nil {
		log.Errf("failed to persist server to config file %s: %s", path, err)
		return
	}

	if written {
		log.Infof("saved server=%s to config file %s", c.server, path)
	} else {
		log.Infof("Warning: config file %s already has a server configured, "+
			"not overriding it with %s", path, c.server)
	}
}
