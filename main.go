package main

import (
	"fmt"
	"os"

	"github.com/railgun-0402/ops-changelog/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
