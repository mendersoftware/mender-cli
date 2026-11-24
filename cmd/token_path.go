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
	"fmt"

	"github.com/spf13/cobra"
)

var tokenPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Print the default on-disk storage path for the authentication token.",
	Long: "Print the absolute filesystem path where mender-cli stores the " +
		"authentication token by default. This always reflects the default " +
		"platform-specific location and ignores any --token override.",
	Run: func(c *cobra.Command, args []string) {
		path, err := getDefaultAuthTokenPath()
		CheckErr(err)
		fmt.Println(path)
	},
}
