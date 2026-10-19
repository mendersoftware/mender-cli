// Copyright 2026 Northern.tech AS
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
	"os"

	"github.com/spf13/cobra"

	"github.com/mendersoftware/mender-cli/client/inventory"
)

var devicesInventoryShowCmd = &cobra.Command{
	Use:   "show DEVICE_ID",
	Short: "Show the inventory of a single device.",
	Args:  cobra.ExactArgs(1),
	Run: func(c *cobra.Command, args []string) {
		cmd, err := NewDevicesInventoryShowCmd(c, args)
		CheckErr(err)
		CheckErr(cmd.Run())
	},
}

func init() {
	addInventoryOutputFlags(devicesInventoryShowCmd)
	addInventoryColumnFlags(devicesInventoryShowCmd)
}

type DevicesInventoryShowCmd struct {
	*inventoryCmdCommon
	deviceID string
}

func NewDevicesInventoryShowCmd(
	cmd *cobra.Command,
	args []string,
) (*DevicesInventoryShowCmd, error) {
	common, err := newInventoryCmdCommon(cmd)
	if err != nil {
		return nil, err
	}
	return &DevicesInventoryShowCmd{inventoryCmdCommon: common, deviceID: args[0]}, nil
}

func (c *DevicesInventoryShowCmd) Run() error {
	client := inventory.NewClient(c.server, c.skipVerify)
	dev, raw, err := client.GetDevice(c.token, c.deviceID)
	if err != nil {
		return err
	}
	if c.rawMode {
		_, err = os.Stdout.Write(append(raw, '\n'))
		return err
	}
	return inventory.Write(os.Stdout, c.output, []inventory.Device{*dev}, c.columns, c.detailLevel)
}
