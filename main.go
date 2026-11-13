package main

import (
	"fmt"
	"os"
)

func main() {
	command := newRootCmd()

	if err := command.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
