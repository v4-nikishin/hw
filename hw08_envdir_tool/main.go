package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Command should be like this '$ go-envdir /path/to/env/dir command arg1 arg2'")
		os.Exit(1)
	}
	m, err := ReadDir(args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if ret := RunCmd(args[2:], m); ret != Success {
		fmt.Println("Failed to run command. Error: ", ret)
	}
}
