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

import "github.com/spf13/cobra"

// tokenCmd is the parent for token-management subcommands.
var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Manage the locally stored authentication token",
	Long: "Manage the locally stored authentication token (e.g. a Personal " +
		"Access Token) without needing to know the on-disk storage location.",
}

func init() {
	tokenCmd.AddCommand(tokenSetCmd)
	tokenCmd.AddCommand(tokenShowCmd)
	tokenCmd.AddCommand(tokenPathCmd)
	tokenCmd.AddCommand(tokenClearCmd)
}
