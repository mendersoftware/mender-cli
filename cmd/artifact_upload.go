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
	"github.com/mendersoftware/mender-cli/log"
)

const (
	argArtifactDescription = "description"
	argWithoutProgress     = "no-progress"
)

var artifactUploadCmd = &cobra.Command{
	Use:   "upload [flags] ARTIFACT",
	Short: "Upload mender artifact to the Mender server.",
	Args:  cobra.ExactArgs(1),
	Run: func(c *cobra.Command, args []string) {
		cmd, err := NewArtifactUploadCmd(c, args)
		CheckErr(err)

		CheckErr(cmd.Run())
	},
}

func init() {
	artifactUploadCmd.Flags().StringP(argArtifactDescription, "", "", "artifact description")
	artifactUploadCmd.Flags().BoolP(argWithoutProgress, "", false, "disable progress bar")
}

type ArtifactUploadCmd struct {
	server          string
	skipVerify      bool
	description     string
	artifactPath    string
	tokenPath       string
	withoutProgress bool
}

func NewArtifactUploadCmd(cmd *cobra.Command, args []string) (*ArtifactUploadCmd, error) {
	server, err := cmd.Flags().GetString(argRootServer)
	if err != nil {
		return nil, err
	}

	skipVerify, err := cmd.Flags().GetBool(argRootSkipVerify)
	if err != nil {
		return nil, err
	}

	artifactDescription, err := cmd.Flags().GetString(argArtifactDescription)
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

	withoutProgress, err := cmd.Flags().GetBool(argWithoutProgress)
	if err != nil {
		return nil, err
	}

	return &ArtifactUploadCmd{
		server:          server,
		description:     artifactDescription,
		tokenPath:       token,
		artifactPath:    args[0],
		skipVerify:      skipVerify,
		withoutProgress: withoutProgress,
	}, nil
}

func (c *ArtifactUploadCmd) Run() error {

	client := deployments.NewClient(c.server, c.skipVerify)
	err := client.UploadArtifact(c.description, c.artifactPath, c.tokenPath, c.withoutProgress)
	if err != nil {
		return err
	}

	log.Info("upload successful")

	return nil
}
