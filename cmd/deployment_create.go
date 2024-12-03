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
	argDeploymentName = "deployment-name"
)

var deploymentCreateCmd = &cobra.Command{
	Use:   "create [flags] DEPLOYMENT_NAME ARTIFACT_NAME DEVICE_LIST",
	Short: "Create a new deployment from an existing artifact on the Mender server.",
	Run: func(c *cobra.Command, args []string) {
		cmd, err := CreateDeploymentCmd(c, args)
		CheckErr(err)
		CheckErr(cmd.Run())
	},
}

func init() {
	deploymentCreateCmd.Flags().StringP(argDeploymentName, "", "",
		"destination path to download to")
}

type DeploymentCreateCmd struct {
	server         string
	deploymentName string
	artifactID     string
	deviceList     string
	token          string
}

func CreateDeploymentCmd(cmd *cobra.Command, args []string) (*DeploymentCreateCmd, error) {
	server := viper.GetString(argRootServer)
	if server == "" {
		return nil, errors.New("No server")
	}

	deploymentName := ""
	if len(args) > 1 {
		deploymentName = args[0]
	}

	artifactID := ""
	if len(args) > 2 {
		artifactID = args[1]
	}

	deviceList := ""
	if len(args) >= 3 {
		deviceList = args[2]
	}

	token, err := getAuthToken(cmd)
	if err != nil {
		return nil, err
	}

	return &DeploymentCreateCmd{
		server:         server,
		deploymentName: deploymentName,
		artifactID:     artifactID,
		deviceList:     deviceList,
		token:          token,
	}, nil
}

func (c *DeploymentCreateCmd) Run() error {
	client := deployments.NewClient(c.server, false)
	err := client.CreateDeployment(c.deploymentName, c.artifactID, c.deviceList, c.token)
	if err != nil {
		return err
	}

	log.Info("deployment created succesfully")

	return nil
}
