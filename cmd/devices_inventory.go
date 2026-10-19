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
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mendersoftware/mender-cli/client/inventory"
)

var devicesInventoryCmd = &cobra.Command{
	Use:       "inventory",
	Short:     "Operations on the inventory of mender devices.",
	ValidArgs: []string{"list", "show", "attributes"},
}

const (
	argOutput  = "output"
	argFilter  = "filter"
	argColumns = "columns"
	argSort    = "sort"
	argAll     = "all"
)

func init() {
	devicesInventoryCmd.AddCommand(devicesInventoryListCmd)
	devicesInventoryCmd.AddCommand(devicesInventoryShowCmd)
	devicesInventoryCmd.AddCommand(devicesInventoryAttributesCmd)
}

func addInventoryOutputFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(argOutput, "o", inventory.FormatText,
		"output format, one of: "+strings.Join(inventory.Formats, ", "))
	cmd.Flags().BoolP(argRawMode, "r", false, "raw mode (json from mender server)")
}

func addInventoryColumnFlags(cmd *cobra.Command) {
	cmd.Flags().IntP(argDetailLevel, "d", 0, "text output detail level [0..1]")
	cmd.Flags().StringSliceP(argColumns, "c", nil,
		"comma-separated attributes to output as [scope/]name; without a scope, "+
			"the first scope having the attribute is used (default: all attributes)")
}

type inventoryCmdCommon struct {
	server      string
	skipVerify  bool
	token       string
	output      string
	rawMode     bool
	detailLevel int
	columns     []inventory.Column
}

func newInventoryCmdCommon(cmd *cobra.Command) (*inventoryCmdCommon, error) {
	server := viper.GetString(argRootServer)
	if server == "" {
		return nil, errors.New("No server")
	}

	flags := cmd.Flags()
	c := &inventoryCmdCommon{server: server}
	var err error

	if c.skipVerify, err = flags.GetBool(argRootSkipVerify); err != nil {
		return nil, err
	}
	if c.output, err = flags.GetString(argOutput); err != nil {
		return nil, err
	}
	if err = inventory.ValidateFormat(c.output); err != nil {
		return nil, err
	}
	if c.rawMode, err = flags.GetBool(argRawMode); err != nil {
		return nil, err
	}
	if flags.Lookup(argDetailLevel) != nil {
		if c.detailLevel, err = flags.GetInt(argDetailLevel); err != nil {
			return nil, err
		}
		if c.detailLevel < 0 || c.detailLevel > 1 {
			return nil, fmt.Errorf("invalid detail level %d, must be 0 or 1", c.detailLevel)
		}
	}
	if flags.Lookup(argColumns) != nil {
		specs, err := flags.GetStringSlice(argColumns)
		if err != nil {
			return nil, err
		}
		if c.columns, err = inventory.ParseColumns(specs); err != nil {
			return nil, err
		}
	}
	if c.token, err = getAuthToken(cmd); err != nil {
		return nil, err
	}
	return c, nil
}

func writeFileAtomically(path string, write func(w io.Writer) error) error {
	f, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+".*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())

	bw := bufio.NewWriter(f)
	err = write(bw)
	if err == nil {
		err = bw.Flush()
	}
	if err == nil {
		err = f.Chmod(0o644)
	}
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	return os.Rename(f.Name(), path)
}
