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
	"github.com/spf13/cobra"

	"github.com/mendersoftware/mender-cli/client/deployments"
)

const (
	argDetailLevel = "detail"
)

var artifactsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List mender artifacts from the Mender server.",
	Run: func(c *cobra.Command, args []string) {
		cmd, err := NewArtifactsListCmd(c, args)
		CheckErr(err)

		CheckErr(cmd.Run())
	},
}

func init() {
	artifactsListCmd.Flags().IntP(argDetailLevel, "d", 0, "artifacts list detail level [0..2]")
}

type ArtifactsListCmd struct {
	server      string
	skipVerify  bool
	tokenPath   string
	detailLevel int
}

func NewArtifactsListCmd(cmd *cobra.Command, args []string) (*ArtifactsListCmd, error) {
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

	detailLevel, err := cmd.Flags().GetInt(argDetailLevel)
	if err != nil {
		return nil, err
	}

	if token == "" {
		token, err = getDefaultAuthTokenPath()
		if err != nil {
			return nil, err
		}
	}

	return &ArtifactsListCmd{
		server:      server,
		tokenPath:   token,
		skipVerify:  skipVerify,
		detailLevel: detailLevel,
	}, nil
}

func (c *ArtifactsListCmd) Run() error {

	client := deployments.NewClient(c.server, c.skipVerify)
	err := client.ListArtifacts(c.tokenPath, c.detailLevel)
	if err != nil {
		return err
	}
	return nil
}
