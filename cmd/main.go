package main

import (
	"fmt"
	"os"
	"ouro/cli" // Replace with your actual package path
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
