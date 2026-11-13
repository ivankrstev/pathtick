package main

import (
	"fmt"
	"os"

	"github.com/ivankrstev/pathtick/internal/cli"
)

func main() {
	command := cli.NewRootCmd()

	if err := command.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
