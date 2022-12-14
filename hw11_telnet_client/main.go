package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", time.Second*10, "connection timeout")
}

func main() {
	flag.Parse()

	var address string

	if len(os.Args) < 3 {
		fmt.Println("Invald args")
		os.Exit(1)
	}

	if len(os.Args) >= 4 {
		address = os.Args[2] + ":" + os.Args[3]
	} else {
		address = os.Args[1] + ":" + os.Args[2]
	}

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		fmt.Printf("Cannot connect: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	fmt.Fprintln(os.Stderr, "...Connected to localhost:4242")
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := client.Receive(); err != nil {
			fmt.Printf("Failed to receive: %v\n", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := client.Send(); err != nil {
			fmt.Printf("Failed to send: %v\n", err)
		}
	}()

	wg.Wait()
	fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
}
