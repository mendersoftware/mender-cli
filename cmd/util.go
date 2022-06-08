// Copyright 2022 Northern.tech AS
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
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func CheckErr(e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "FAILURE: %s\n", e.Error())
		os.Exit(1)
	}
}

func migrateAuthToken(oldtoken string, token string) {
	// if needed, migrate token from old to new location
	if _, err := os.Stat(token); !os.IsNotExist(err) {
		// new token exists, no migration
		return
	}

	if _, err := os.Stat(oldtoken); err != nil {
		// old token doesn't exist, no migration
		return
	}

	// Attempt migration, ignore errors (but log them?)
	if err := os.MkdirAll(filepath.Dir(token), 0700); err == nil {
		// log that token was moved?
		_ = os.Rename(oldtoken, token)
	}

	// Cleanup old token directory if empty
	_ = os.Remove(filepath.Dir(oldtoken)) // err on non-empty, ignore.
}

func getDefaultAuthTokenPath() (string, error) {
	cachedir := ""
	userhomedir := ""

	if homeenv := os.Getenv("HOME"); homeenv != "" {
		userhomedir = homeenv
	} else if user, err := user.Current(); err == nil {
		userhomedir = user.HomeDir
	} else {
		return "", errors.New("Not able to determine users cache dir")
	}

	if cachehomeenv := os.Getenv("XDG_CACHE_HOME"); cachehomeenv != "" {
		cachedir = cachehomeenv
	} else {
		cachedir = path.Join(userhomedir, ".cache")
	}

	oldtoken := filepath.Join(userhomedir, ".mender", "authtoken")
	token := filepath.Join(cachedir, "mender", "authtoken")

	migrateAuthToken(oldtoken, token)

	return token, nil
}

func getAuthToken(cmd *cobra.Command) (string, error) {
	tokenValue, err := cmd.Flags().GetString(argRootTokenValue)
	if err != nil {
		return "", err
	} else if tokenValue != "" {
		return tokenValue, nil
	}
	tokenPath, err := getDefaultAuthTokenPath()
	if err != nil {
		return "", err
	}
	token, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return "", errors.Wrap(err, "Please Login first")
	}
	return string(token), nil
}
