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
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/apcera/termtables"
	"github.com/spf13/cobra"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// TODO: add line number
type hbaRule struct {
	connectionType string
	databaseName   string
	userName       string
	ipAddress      string
	networkMask    string
	authType       string
	lineNumber     int
}

type hbaRules []hbaRule

// more details here: http://nerdyworm.com/blog/2013/05/15/sorting-a-slice-of-structs-in-go/

func (slice hbaRules) Len() int {
	return len(slice)
}

func (slice hbaRules) Less(i, j int) bool {
	return slice[i].connectionType < slice[j].connectionType
}

func (slice hbaRules) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func openFile(filename string) hbaRules {
	file, err := os.Open(filename)
	check(err)
	defer file.Close()

	fileRules := []hbaRule{}

	scanner := bufio.NewScanner(file)

	lineNo := 0

	for scanner.Scan() {
		newLine := scanner.Text()
		lineNo += 1

		if strings.HasPrefix(newLine, "host") || strings.HasPrefix(newLine, "local") {

			re := regexp.MustCompile(`[A-Za-z0-9_'\./\+\:]+`)
			matches := re.FindAllString(newLine, -1)

			newRule := hbaRule{
				connectionType: matches[0],
				databaseName:   matches[1],
				userName:       matches[2],
				lineNumber:     lineNo}

			if len(matches) == 6 { // full mask
				newRule.ipAddress = matches[3]
				newRule.networkMask = matches[4]
				newRule.authType = matches[5]
			} else if len(matches) == 5 { // /32 type mask

				s := strings.Split(matches[3], "/")
				newRule.ipAddress = s[0]
				newRule.networkMask = "/" + s[1]
				newRule.authType = matches[4]
			} else { // local connection
				newRule.authType = matches[3]
			}

			fileRules = append(fileRules, newRule)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return fileRules
}

func printSlice(s []hbaRule) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Find a rule on the pg_hba.conf file.",
	Long:  `Find a rule on the pg_hba.conf file.`,
	Run: func(cmd *cobra.Command, args []string) {
		hba_file := pgdata + "/pg_hba.conf"

		// TODO: Work your own magic here
		// fmt.Println("Using the hba file: ", hba_file)

		rules := openFile(hba_file)
		sort.Sort(rules)
		table := termtables.CreateTable()

		table.AddHeaders("Line", "Type", "Database", "User/Group", "Host", "Mask", "Method")

		for _, element := range rules {
			table.AddRow(element.lineNumber, element.connectionType, element.databaseName, element.userName, element.ipAddress, element.networkMask, element.authType)
		}
		fmt.Println(table.Render())

		// local: (?P<type>(local))\s+(?P<database>(\w+))\s+(?P<user>(\w+))\s+(?P<mode>(\w+))

	},
}

func init() {
	RootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
