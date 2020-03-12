// Copyright 2020 Northern.tech AS
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
	"path/filepath"

	"github.com/howeyc/gopass"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/mendersoftware/mender-cli/client/useradm"
	"github.com/mendersoftware/mender-cli/log"
)

const (
	argLoginUsername = "username"
	argLoginPassword = "password"
	argLoginToken    = "token"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to the Mender server (required before other operations).",
	Run: func(c *cobra.Command, args []string) {
		cmd, err := NewLoginCmd(c, args)
		CheckErr(err)

		CheckErr(cmd.Run())
	},
}

func init() {
	loginCmd.Flags().StringP(argLoginUsername, "", "", "username, format: email (required)")
	loginCmd.MarkFlagRequired(argLoginUsername)

	loginCmd.Flags().StringP(argLoginPassword, "", "", "password (will prompt if not provided)")
	loginCmd.Flags().StringP(argLoginToken, "", "", "two-factor authentication token")
}

type LoginCmd struct {
	server     string
	skipVerify bool
	username   string
	password   string
	token      string
	tokenPath  string
}

func NewLoginCmd(cmd *cobra.Command, args []string) (*LoginCmd, error) {
	server, err := cmd.Flags().GetString(argRootServer)
	if err != nil {
		return nil, err
	}

	skipVerify, err := cmd.Flags().GetBool(argRootSkipVerify)
	if err != nil {
		return nil, err
	}

	username, err := cmd.Flags().GetString(argLoginUsername)
	if err != nil {
		return nil, err
	}

	password, err := cmd.Flags().GetString(argLoginPassword)
	if err != nil {
		return nil, err
	}

	tfaToken, err := cmd.Flags().GetString(argLoginToken)
	if err != nil {
		return nil, err
	}

	token, err := cmd.Flags().GetString(argRootToken)
	if err != nil {
		return nil, err
	}

	if token == "" {
		token, err = getDefaultAuthTokenPath()
		if err != nil {
			return nil, err
		}
	}

	return &LoginCmd{
		server:     server,
		username:   username,
		password:   password,
		token:      tfaToken,
		tokenPath:  token,
		skipVerify: skipVerify,
	}, nil
}

func (c *LoginCmd) Run() error {
	err := c.maybeGetPassword()
	if err != nil {
		return err
	}

	client := useradm.NewClient(c.server, c.skipVerify)
	res, err := client.Login(c.username, c.password, c.token)
	if err != nil {
		return err
	}

	err = c.saveToken(res)
	if err != nil {
		return err
	}

	return nil
}

func (c *LoginCmd) maybeGetPassword() error {
	if c.password == "" {
		fmt.Printf("Password:")

		p, err := gopass.GetPasswdMasked()
		if err != nil {
			return err
		}

		c.password = string(p)
	}

	return nil
}

func (c *LoginCmd) saveToken(t []byte) error {
	dir := filepath.Dir(c.tokenPath)
	log.Verbf("creating directory: %v\n", dir)

	err := os.MkdirAll(dir, os.ModeDir|0700)
	if err != nil {
		return errors.Wrapf(err, "failed to create directory %s", dir)

	}

	err = ioutil.WriteFile(c.tokenPath, t, 0600)
	if err != nil {
		return errors.Wrapf(err, "failed to create file %s", c.tokenPath)
	}

	log.Verb("saved token to: " + c.tokenPath)
	log.Info("login successful")

	return nil
}
