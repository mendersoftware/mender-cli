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

	"github.com/mendersoftware/mender-cli/log"
)

const (
	argRootServer     = "server"
	argRootSkipVerify = "skip-verify"
	argRootToken      = "token"
	argRootVerbose    = "verbose"

	defaultTokenPath = "/tmp/mendersoftware/authtoken"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mender-cli",
	Short: "A general-purpose CLI for the Mender server.",

	//setup global stuff, will run regardless of (sub)command
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		verbose, err := cmd.Flags().GetBool(argRootVerbose)
		CheckErr(err)
		log.Setup(verbose)

		if verbose {
			log.Verb(fmt.Sprintf("verbose output is ON"))
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.AddCommand(loginCmd)
	CheckErr(rootCmd.Execute())
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringP(argRootServer, "", "", "root server URL, e.g. 'https://hosted.mender.io' (required)")
	rootCmd.MarkPersistentFlagRequired(argRootServer)
	rootCmd.PersistentFlags().BoolP(argRootSkipVerify, "k", false, "skip SSL certificate verification")
	rootCmd.PersistentFlags().StringP(argRootToken, "", "", "token file path")
	rootCmd.PersistentFlags().BoolP(argRootVerbose, "v", false, "print verbose output")
}
