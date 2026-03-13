package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ops-changelog",
	Short: "Quickly see what changed in a service — useful during incidents",
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
