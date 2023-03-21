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
	argArtifactDescription = "description"
	argWithoutProgress     = "no-progress"
	argDirect              = "direct"
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
	artifactUploadCmd.Flags().BoolP(argDirect, "", false, "upload directly to storage")
}

type ArtifactUploadCmd struct {
	server          string
	skipVerify      bool
	description     string
	artifactPath    string
	token           string
	withoutProgress bool
	direct          bool
}

func NewArtifactUploadCmd(cmd *cobra.Command, args []string) (*ArtifactUploadCmd, error) {
	server := viper.GetString(argRootServer)
	if server == "" {
		return nil, errors.New("No server")
	}

	skipVerify, err := cmd.Flags().GetBool(argRootSkipVerify)
	if err != nil {
		return nil, err
	}

	artifactDescription, err := cmd.Flags().GetString(argArtifactDescription)
	if err != nil {
		return nil, err
	}

	withoutProgress, err := cmd.Flags().GetBool(argWithoutProgress)
	if err != nil {
		return nil, err
	}

	direct, err := cmd.Flags().GetBool(argDirect)
	if err != nil {
		return nil, err
	}

	token, err := getAuthToken(cmd)
	if err != nil {
		return nil, err
	}

	return &ArtifactUploadCmd{
		server:          server,
		description:     artifactDescription,
		token:           token,
		artifactPath:    args[0],
		skipVerify:      skipVerify,
		withoutProgress: withoutProgress,
		direct:          direct,
	}, nil
}

func (c *ArtifactUploadCmd) Run() error {
	client := deployments.NewClient(c.server, c.skipVerify)
	if c.direct {
		log.Infof("getting direct link.\n")
		link, err := client.DirectDownloadLink(c.token)
		if err != nil {
			return errors.Wrap(err, "failed to get the direct pre-signed URL")
		}

		log.Infof("uploading the artifact.\n")
		err = client.DirectUpload(
			c.token,
			c.artifactPath,
			link.ArtifactID,
			link.Uri,
			link.Header,
			c.withoutProgress,
		)
		if err != nil {
			return errors.Wrap(err, "failed to upload the artifact")
		}
	} else {
		err := client.UploadArtifact(c.description, c.artifactPath, c.token, c.withoutProgress)
		if err != nil {
			return err
		}
	}

	log.Info("upload successful")

	return nil
}
