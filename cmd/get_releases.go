// Copyright 2018 Northern.tech AS
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
	"bytes"
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/mendersoftware/mender-cli/client/deployments"
	"github.com/mendersoftware/mender-cli/log"
)

var getReleasesCmd = &cobra.Command{
	Use:   "get",
	Short: "Get list of mender releases.",
	Run: func(c *cobra.Command, args []string) {
		cmd, err := NewGetReleasesCmd(c, args)
		CheckErr(err)

		CheckErr(cmd.Run())
	},
}

type GetReleasesCmd struct {
	server     string
	skipVerify bool
	tokenPath  string
}

func NewGetReleasesCmd(cmd *cobra.Command, args []string) (*GetReleasesCmd, error) {
	server, err := cmd.Flags().GetString(argRootServer)
	if err != nil {
		return nil, err
	}

	skipVerify, err := cmd.Flags().GetBool(argRootSkipVerify)
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

	return &GetReleasesCmd{
		server:     server,
		tokenPath:  token,
		skipVerify: skipVerify,
	}, nil
}

func (c *GetReleasesCmd) Run() error {
	client := deployments.NewClient(c.server, c.skipVerify)
	releases, err := client.GetReleases(c.tokenPath)
	if err != nil {
		return err
	}
	var releasesJson bytes.Buffer
	if err := json.Indent(&releasesJson, releases, "", "  "); err != nil {
		return err
	}
	log.Info("Releases:")
	log.Info(string(releasesJson.Bytes()))

	return nil
}
