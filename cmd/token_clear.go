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
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mendersoftware/mender-cli/log"
)

const argTokenClearYes = "yes"

var tokenClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Delete the locally stored authentication token.",
	Long: "Delete the locally stored authentication token. By default an " +
		"interactive confirmation is required; use --yes to skip it. In " +
		"non-interactive contexts --yes is mandatory.",
	Run: func(c *cobra.Command, args []string) {
		yes, err := c.Flags().GetBool(argTokenClearYes)
		CheckErr(err)
		CheckErr(runTokenClear(c, yes))
	},
}

func init() {
	tokenClearCmd.Flags().BoolP(argTokenClearYes, "y", false,
		"do not prompt for confirmation before deleting")
}

func runTokenClear(cmd *cobra.Command, yes bool) error {
	path, err := resolveTokenPath(cmd)
	if err != nil {
		return err
	}
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			log.Info("no token stored")
			return nil
		}
		return fmt.Errorf("failed to stat token file: %w", err)
	}

	if !yes {
		if !isStdinTTY() {
			return errors.New(
				"refusing to delete token without --yes in non-interactive mode")
		}
		fmt.Fprintf(os.Stderr, "Delete stored token at %s? [y/N]: ", path)
		reader := bufio.NewReader(os.Stdin)
		answer, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read confirmation: %w", err)
		}
		answer = strings.ToLower(strings.TrimSpace(answer))
		if answer != "y" && answer != "yes" {
			log.Info("aborted; token not deleted")
			return nil
		}
	}

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}
	log.Info("token cleared")
	return nil
}
