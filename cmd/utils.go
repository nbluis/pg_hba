package cmd

import "fmt"

func verboseLog(a ...interface{}) {
	if verboseMode {
		fmt.Println(a...)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
