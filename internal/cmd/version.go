package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Build   string
	Version string
)

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Show build and version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Build: %s\nVersion: %s\n", Build, Version)
	},
}
