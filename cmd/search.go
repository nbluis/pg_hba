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
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/apcera/termtables"
	"github.com/spf13/cobra"
)

// TODO: add line number
type hbaRule struct {
	connectionType string
	databaseName   string
	userName       string
	ipAddress      string
	networkMask    string
	authType       string
	lineNumber     int
	comments       string
}

// more details about sort : http://nerdyworm.com/blog/2013/05/15/sorting-a-slice-of-structs-in-go/
//     and: http://stackoverflow.com/questions/28999735/what-is-the-shortest-way-to-simply-sort-an-array-of-structs-by-arbitrary-field

func processFile(filename string) ([]hbaRule, error) {
	file, err := os.Open(filename)
	defer file.Close()

	fileRules := []hbaRule{}

	if err != nil {
		return fileRules, err
	}

	scanner := bufio.NewScanner(file)

	lineNo := 0

	for scanner.Scan() {
		newLine := scanner.Text()
		lineNo += 1

		if strings.HasPrefix(newLine, "host") || strings.HasPrefix(newLine, "local") {

			// re := regexp.MustCompile(`[A-Za-z0-9_'\./\+\:\,]+`)
			// matches := re.FindAllString(newLine, -1)
			comments := strings.Split(newLine, "#")

			// if len(comments) > 1 {
			matches := strings.Fields(comments[0])
			// } else

			// fmt.Printf("len=%d cap=%d %v\n", len(comments), cap(comments), comments)

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

			if len(comments) > 1 {
				if !foundComments {
					foundComments = true
				}
				newRule.comments = comments[1]
			}

			fileRules = append(fileRules, newRule)

			// limit the max results
			if rowLimit > 0 && len(fileRules) >= rowLimit {
				break
			}
		} else {
			err := errors.New("emit macho dwarf: elf header corrupted")
			return fileRules, err
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return fileRules, nil
}

var (
	foundComments bool = false
)

func removeComments(content []byte) []byte {
	cppcmt := regexp.MustCompile(`#.*`)
	return cppcmt.ReplaceAll(content, []byte(""))
}

func printSlice(s []hbaRule) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}

func formatRules(ruleList []hbaRule) string {

	table := termtables.CreateTable()

	if foundComments {
		table.AddHeaders("Line", "Type", "Database", "User/Group", "Host", "Mask", "Method", "Comment")
	} else {
		table.AddHeaders("Line", "Type", "Database", "User/Group", "Host", "Mask", "Method")
	}

	for _, element := range ruleList {
		if foundComments {
			table.AddRow(element.lineNumber, element.connectionType, element.databaseName, element.userName, element.ipAddress, element.networkMask, element.authType, element.comments)
		} else {
			table.AddRow(element.lineNumber, element.connectionType, element.databaseName, element.userName, element.ipAddress, element.networkMask, element.authType)
		}
	}
	return table.Render()
}

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Find a rule on the pg_hba.conf file.",
	Long: `Find a rule on the pg_hba.conf file.

Valid sorting options:
 - connectionType
 - databaseName
 - userName
 - ipAddress
 - networkMask
 - authType
 - lineNumber
 - comments

	`,
	Run: func(cmd *cobra.Command, args []string) {
		hba_file := pgdata + "/pg_hba.conf"

		verboseLog("Using the hba file: ", hba_file)

		fileRules, err := processFile(hba_file)

		check(err)

		rules := sortRules(fileRules)

		fmt.Println(formatRules(rules))

		verboseLog("("+strconv.Itoa(len(rules)), "rows found)")

		// local: (?P<type>(local))\s+(?P<database>(\w+))\s+(?P<user>(\w+))\s+(?P<mode>(\w+))

	},
}

func sortRules(ruleList []hbaRule) []hbaRule {
	switch sortOrder {
	case "connectionType":
		sort.Sort(connectionTypeSorter(ruleList))
	case "databaseName":
		sort.Sort(databaseNameSorter(ruleList))
	case "userName":
		sort.Sort(userNameSorter(ruleList))
	case "ipAddress":
		sort.Sort(ipAddressSorter(ruleList))
	case "networkMask":
		sort.Sort(networkMaskSorter(ruleList))
	case "authType":
		sort.Sort(authTypeSorter(ruleList))
	case "lineNumber":
		sort.Sort(lineNumberSorter(ruleList))
	case "comments":
		sort.Sort(commentsSorter(ruleList))
	default:
		sort.Sort(lineNumberSorter(ruleList))
	}

	return ruleList
}

var (
	sortOrder string = "lineNumber"
	rowLimit  int    = -1
)

func init() {
	RootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringVarP(&sortOrder, "sort-by", "s", sortOrder, "Change the sort order")
	searchCmd.Flags().IntVarP(&rowLimit, "limit", "l", rowLimit, "Limit the result for x rows")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
