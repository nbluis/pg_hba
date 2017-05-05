// Copyright © 2017 Sebastian Webber <sebastian@swebber.me>
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
//    this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors
//    may be used to endorse or promote products derived from this software
//    without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [description]",
	Short: "Add a new rule on the pg_hba.conf file",
	Long:  `Add a new rule on the pg_hba.conf file`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Missing the description.")
		} else {

			newData := hbaRule{
				authType:     "host",
				userName:     username,
				databaseName: database,
				ipAddress:    hostAddress,
				comments:     args[0]}

			fmt.Println(newData)
			// fmt.Printf("## %s \n", args[0])
			// fmt.Printf("host   %s    %s    %s/32   md5\n", username, database, hostAddress)
		}
	},
}

var (
	hostAddress = "127.0.0.1"
	username    = "user"
	database    = "all"
)

func init() {
	RootCmd.AddCommand(addCmd)

	// addCmd.SetArgs([]string{"sub", "arg1", "arg2"})

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("username", "", "A help for foo")
	// addCmd.ar
	addCmd.Flags().StringVarP(&hostAddress, "host", "h", hostAddress, "Host Address")
	addCmd.Flags().StringVarP(&username, "username", "U", username, "Username")
	addCmd.Flags().StringVarP(&database, "database", "d", database, "Database name")

	addCmd.Flags().BoolP("help", "H", false, "Help message add")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:

}
