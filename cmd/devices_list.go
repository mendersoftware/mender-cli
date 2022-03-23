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
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mendersoftware/mender-cli/client/devices"
)

var devicesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Get a list of devices from the Mender server.",
	Run: func(c *cobra.Command, args []string) {
		cmd, err := NewDevicesListCmd(c, args)
		CheckErr(err)
		CheckErr(cmd.Run())
	},
}

const argRawMode = "raw"

func init() {
	devicesListCmd.Flags().IntP(argDetailLevel, "d", 0, "devices list detail level [0..3]")
	devicesListCmd.Flags().BoolP(
		argRawMode,
		"r",
		false,
		"devices list raw mode (json from mender server)")
}

type DevicesListCmd struct {
	server      string
	skipVerify  bool
	tokenPath   string
	detailLevel int
	rawMode     bool
}

func NewDevicesListCmd(cmd *cobra.Command, args []string) (*DevicesListCmd, error) {
	server := viper.GetString(argRootServer)
	if server == "" {
		return nil, errors.New("No server")
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

	rawMode, err := cmd.Flags().GetBool(argRawMode)
	if err != nil {
		return nil, err
	}

	if token == "" {
		token, err = getDefaultAuthTokenPath()
		if err != nil {
			return nil, err
		}
	}

	return &DevicesListCmd{
		server:      server,
		tokenPath:   token,
		skipVerify:  skipVerify,
		detailLevel: detailLevel,
		rawMode:     rawMode,
	}, nil
}

func (c *DevicesListCmd) Run() error {

	client := devices.NewClient(c.server, c.skipVerify)
	return client.ListDevices(c.tokenPath, c.detailLevel, c.rawMode)
}
