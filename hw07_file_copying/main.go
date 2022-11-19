package main

import (
	"flag"
	"fmt"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	fmt.Printf("Copying %s to %s\n", from, to)

	if from == "" || to == "" {
		fmt.Println("Input parameter error, usage method: the executable program copies the source file to the target file")
		flag.CommandLine.Usage()
		return
	}

	if from == to {
		fmt.Println("The name of the target file and the name of the source file cannot match")
		flag.CommandLine.Usage()
		return
	}

	err := Copy(from, to, offset, limit)
	if err != nil {
		fmt.Printf("Copyimg failed %q\n", err)
	} else {
		fmt.Printf("Copying succeeded\n")
	}
}
