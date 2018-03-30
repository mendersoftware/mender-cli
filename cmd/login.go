// Copyright 2017 Northern.tech AS
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
	"fmt"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to the Mender backend (required before other operations).",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("login called")
		cmd.Flags().GetString("server")
		cmd.Flags().GetString("username")
	},
}

func init() {
	loginCmd.Flags().StringP("username", "", "", "username, format: email (required)")
	loginCmd.MarkFlagRequired("username")

	loginCmd.Flags().StringP("password", "", "", "password")
}
