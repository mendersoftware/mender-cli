// Copyright 2023 Northern.tech AS
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
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mendersoftware/mender-cli/client/deployments"
	"github.com/mendersoftware/mender-cli/log"
)

var artifactDeleteCmd = &cobra.Command{
	Use:   "delete [flags] ARTIFACT_ID",
	Short: "Delete mender artifact from the Mender server.",
	Args:  cobra.ExactArgs(1),
	Run: func(c *cobra.Command, args []string) {
		cmd, err := NewArtifactDeleteCmd(c, args)
		CheckErr(err)
		CheckErr(cmd.Run())
	},
}

func init() {
}

type ArtifactDeleteCmd struct {
	server     string
	skipVerify bool
	artifactID string
	token      string
}

func NewArtifactDeleteCmd(cmd *cobra.Command, args []string) (*ArtifactDeleteCmd, error) {
	server := viper.GetString(argRootServer)
	if server == "" {
		return nil, errors.New("No server")
	}

	skipVerify, err := cmd.Flags().GetBool(argRootSkipVerify)
	if err != nil {
		return nil, err
	}

	token, err := getAuthToken(cmd)
	if err != nil {
		return nil, err
	}

	artifactID := ""
	if len(args) == 1 {
		artifactID = args[0]
	}

	return &ArtifactDeleteCmd{
		server:     server,
		artifactID: artifactID,
		token:      token,
		skipVerify: skipVerify,
	}, nil
}

func (c *ArtifactDeleteCmd) Run() error {

	client := deployments.NewClient(c.server, c.skipVerify)
	err := client.DeleteArtifact(c.artifactID, c.token)
	if err != nil {
		return err
	}

	log.Info("delete successful")

	return nil
}
