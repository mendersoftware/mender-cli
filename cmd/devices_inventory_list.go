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
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/mendersoftware/mender-cli/client/inventory"
)

const (
	allPerPage     = 500
	exitIncomplete = 3

	argProgress   = "progress"
	argCount      = "count"
	argOutputFile = "output-file"
)

var stdoutIsTerminal = func() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

var devicesInventoryListCmd = &cobra.Command{
	Use:   "list",
	Short: "Search the inventory of devices on the Mender server.",
	Long: `Search the inventory of devices on the Mender server.

Filters are combined with AND; an attribute without a scope is looked up in
the inventory scope. Supported filter expressions:

  [scope/]name=value       equal
  [scope/]name!=value      not equal
  [scope/]name=~regex      regular expression (anchored at the start)
  [scope/]name>value       greater than (also >=, <, <=)
  [scope/]name=in:a,b,c    equal to any of the values
  [scope/]name=nin:a,b,c   equal to none of the values
  [scope/]name?            attribute exists
  ![scope/]name?           attribute does not exist

Scopes are: inventory, identity, system, tags and monitor. Values are compared
as strings. Array attribute values are joined with ';' in text and CSV output.

Devices are written as each page arrives, except for JSON, CSV and TSV output
without --columns, which is written once all devices are fetched.

With --all, the exit code is 3 if fewer devices were fetched than the server
reported, because the inventory changed while fetching; the output is still
written.`,
	Example: `  mender-cli devices inventory list \
    --filter device_type=<device type> --filter identity/status=accepted \
    --columns hostname,artifact_name,system/check_in_time,identity/mac \
    --sort artifact_name --all --output csv

  mender-cli devices inventory list --all --progress --output csv \
    --filter device_type=<device type> --columns hostname,artifact_name \
    --output-file devices.csv

  mender-cli devices inventory list --count --filter device_type=<device type>`,
	Run: func(c *cobra.Command, args []string) {
		cmd, err := NewDevicesInventoryListCmd(c, args)
		CheckErr(err)
		CheckErr(cmd.Run())
	},
}

func init() {
	addInventoryOutputFlags(devicesInventoryListCmd)
	addInventoryColumnFlags(devicesInventoryListCmd)
	flags := devicesInventoryListCmd.Flags()
	flags.StringArrayP(argFilter, "f", nil,
		"filter expression, can be repeated (see above)")
	flags.StringSliceP(argSort, "s", nil,
		"comma-separated attributes to sort by as [scope/]name[:asc|:desc]")
	flags.IntP(argPerPage, "N", 20, "Number of results to display")
	flags.IntP(argPage, "P", 1, "Page number to return")
	flags.BoolP(argAll, "a", false,
		fmt.Sprintf("fetch all pages (%d devices per page unless --%s is given)",
			allPerPage, argPerPage))
	flags.Bool(argProgress, false,
		fmt.Sprintf("print the percentage of fetched devices to stderr after every page; "+
			"requires --%s and --%s or redirected stdout", argAll, argOutputFile))
	flags.Bool(argCount, false, "only print the number of matching devices")
	flags.String(argOutputFile, "",
		"write the output to this file, replacing it only if fetching succeeds")
	for _, f := range []string{argAll, argProgress, argOutput, argColumns, argRawMode} {
		devicesInventoryListCmd.MarkFlagsMutuallyExclusive(argCount, f)
	}
}

type DevicesInventoryListCmd struct {
	*inventoryCmdCommon
	params     inventory.SearchParams
	all        bool
	progress   bool
	count      bool
	outputFile string
	stdout     io.Writer
	stderr     io.Writer
}

func NewDevicesInventoryListCmd(
	cmd *cobra.Command,
	args []string,
) (*DevicesInventoryListCmd, error) {
	common, err := newInventoryCmdCommon(cmd)
	if err != nil {
		return nil, err
	}
	c := &DevicesInventoryListCmd{
		inventoryCmdCommon: common,
		stdout:             os.Stdout,
		stderr:             os.Stderr,
	}
	flags := cmd.Flags()

	if c.params.PerPage, err = flags.GetInt(argPerPage); err != nil {
		return nil, err
	}
	if c.params.Page, err = flags.GetInt(argPage); err != nil {
		return nil, err
	}
	if c.all, err = flags.GetBool(argAll); err != nil {
		return nil, err
	}
	if c.all && !flags.Changed(argPerPage) {
		c.params.PerPage = allPerPage
	}
	if c.params.Page <= 0 || c.params.PerPage <= 0 {
		return nil, errors.New("page and per-page arguments must be larger than 0")
	}
	if c.progress, err = flags.GetBool(argProgress); err != nil {
		return nil, err
	}
	if c.count, err = flags.GetBool(argCount); err != nil {
		return nil, err
	}
	if c.outputFile, err = flags.GetString(argOutputFile); err != nil {
		return nil, err
	}
	if err = validateListOptions(c.all, c.progress, c.outputFile, stdoutIsTerminal()); err != nil {
		return nil, err
	}

	filters, err := flags.GetStringArray(argFilter)
	if err != nil {
		return nil, err
	}
	for _, f := range filters {
		pred, err := inventory.ParseFilter(f)
		if err != nil {
			return nil, err
		}
		c.params.Filters = append(c.params.Filters, pred)
	}

	sort, err := flags.GetStringSlice(argSort)
	if err != nil {
		return nil, err
	}
	if c.params.Sort, err = inventory.ParseSort(sort); err != nil {
		return nil, err
	}
	c.params.Sort = inventory.StableSort(c.params.Sort)

	c.params.Attributes = inventory.SelectAttributes(c.columns, c.params.Sort)
	return c, nil
}

func validateListOptions(all, progress bool, outputFile string, stdoutIsTerminal bool) error {
	if progress && !all {
		return fmt.Errorf("--%s requires --%s", argProgress, argAll)
	}
	if progress && outputFile == "" && stdoutIsTerminal {
		return fmt.Errorf("--%s needs --%s or stdout redirected to a file or pipe",
			argProgress, argOutputFile)
	}
	return nil
}

func progressLine(fetched, total int) string {
	percent := 100
	if total > 0 && fetched < total {
		percent = fetched * 100 / total
	}
	return fmt.Sprintf("Fetched %d/%d devices (%d%%)", fetched, total, percent)
}

func (c *DevicesInventoryListCmd) Run() error {
	var incomplete *exitError
	write := func(w io.Writer) error {
		var err error
		incomplete, err = c.write(w)
		return err
	}
	var err error
	if c.outputFile != "" {
		err = writeFileAtomically(c.outputFile, write)
	} else {
		err = write(c.stdout)
	}
	if err != nil {
		return err
	}
	if incomplete != nil {
		return incomplete
	}
	return nil
}

func (c *DevicesInventoryListCmd) write(w io.Writer) (*exitError, error) {
	client := inventory.NewClient(c.server, c.skipVerify)

	if c.count {
		params := c.params
		params.Page, params.PerPage = 1, 1
		res, err := client.Search(c.token, params)
		if err != nil {
			return nil, err
		}
		_, err = fmt.Fprintln(w, res.Total)
		return nil, err
	}

	wr, err := inventory.NewWriter(w, c.output, c.columns, c.detailLevel, c.rawMode)
	if err != nil {
		return nil, err
	}

	if !c.all {
		res, err := client.Search(c.token, c.params)
		if err != nil {
			return nil, err
		}
		if err := wr.WritePage(res.Devices, res.Raw); err != nil {
			return nil, err
		}
		return nil, wr.Close()
	}

	skipped := (c.params.Page - 1) * c.params.PerPage
	stats, err := client.SearchEach(c.token, c.params,
		func(page *inventory.SearchResult, stats inventory.SearchStats) error {
			if err := wr.WritePage(page.Devices, page.Raw); err != nil {
				return err
			}
			if c.progress {
				fmt.Fprintln(c.stderr, progressLine(stats.Fetched, stats.Total-skipped))
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	if err := wr.Close(); err != nil {
		return nil, err
	}
	if expected := stats.Total - skipped; stats.Fetched < expected {
		msg := fmt.Sprintf("fetched %d of %d devices, the inventory changed while fetching",
			stats.Fetched, expected)
		if stats.Duplicates > 0 {
			msg = fmt.Sprintf("fetched %d of %d devices, %d were returned more than once; "+
				"the server returned inconsistent pages or the inventory changed while fetching",
				stats.Fetched, expected, stats.Duplicates)
		}
		return &exitError{
			err:   errors.New(msg),
			code:  exitIncomplete,
			label: "WARNING",
		}, nil
	}
	return nil, nil
}
