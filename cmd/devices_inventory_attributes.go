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

var devicesInventoryAttributesCmd = &cobra.Command{
	Use:   "attributes",
	Short: "List the inventory attributes which can be used in filters and columns.",
	Long: `List the inventory attributes which can be used in filters and columns,
most common first.

The Mender server samples at most 5000 devices and returns at most 500
attributes.`,
	Run: func(c *cobra.Command, args []string) {
		cmd, err := NewDevicesInventoryAttributesCmd(c, args)
		CheckErr(err)
		CheckErr(cmd.Run())
	},
}

func init() {
	addInventoryOutputFlags(devicesInventoryAttributesCmd)
}

type DevicesInventoryAttributesCmd struct {
	*inventoryCmdCommon
}

func NewDevicesInventoryAttributesCmd(
	cmd *cobra.Command,
	args []string,
) (*DevicesInventoryAttributesCmd, error) {
	common, err := newInventoryCmdCommon(cmd)
	if err != nil {
		return nil, err
	}
	return &DevicesInventoryAttributesCmd{inventoryCmdCommon: common}, nil
}

func (c *DevicesInventoryAttributesCmd) Run() error {
	client := inventory.NewClient(c.server, c.skipVerify)
	attrs, raw, err := client.FilterableAttributes(c.token)
	if err != nil {
		return err
	}
	if c.rawMode {
		_, err = os.Stdout.Write(append(raw, '\n'))
		return err
	}
	return inventory.WriteFilterableAttributes(os.Stdout, c.output, attrs)
}
