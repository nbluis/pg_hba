package cmd

import "fmt"

func verboseLog(a ...interface{}) {
	if verboseMode {
		fmt.Println(a...)
	}
}
