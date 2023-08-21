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
	"github.com/pkg/errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mendersoftware/mender-cli/client/deployments"
	"github.com/mendersoftware/mender-cli/log"
)

const (
	argDestinationPath = "destination-path"
)

var artifactDownloadCmd = &cobra.Command{
	Use:   "download [flags] ARTIFACT",
	Short: "Download mender artifact from the Mender server.",
	Run: func(c *cobra.Command, args []string) {
		cmd, err := NewArtifactDownloadCmd(c, args)
		CheckErr(err)
		CheckErr(cmd.Run())
	},
}

func init() {
	artifactDownloadCmd.Flags().StringP(argDestinationPath, "", "",
		"destination path to download to")
	artifactDownloadCmd.Flags().BoolP(argWithoutProgress, "", false,
		"disable progress bar")
}

type ArtifactDownloadCmd struct {
	server          string
	skipVerify      bool
	destinationPath string
	artifactID      string
	token           string
	withoutProgress bool
}

func NewArtifactDownloadCmd(cmd *cobra.Command, args []string) (*ArtifactDownloadCmd, error) {
	server := viper.GetString(argRootServer)
	if server == "" {
		return nil, errors.New("No server")
	}

	skipVerify, err := cmd.Flags().GetBool(argRootSkipVerify)
	if err != nil {
		return nil, err
	}

	withoutProgress, err := cmd.Flags().GetBool(argWithoutProgress)
	if err != nil {
		return nil, err
	}

	destinationPath, err := cmd.Flags().GetString(argDestinationPath)
	if err != nil {
		return nil, err
	}

	artifactID := ""
	if len(args) == 1 {
		artifactID = args[0]
	}

	token, err := getAuthToken(cmd)
	if err != nil {
		return nil, err
	}

	return &ArtifactDownloadCmd{
		server:          server,
		destinationPath: destinationPath,
		artifactID:      artifactID,
		token:           token,
		skipVerify:      skipVerify,
		withoutProgress: withoutProgress,
	}, nil
}

func (c *ArtifactDownloadCmd) Run() error {
	client := deployments.NewClient(c.server, c.skipVerify)
	err := client.DownloadArtifact(c.destinationPath, c.artifactID, c.token, c.withoutProgress)
	if err != nil {
		return err
	}

	log.Info("download successful")

	return nil
}
