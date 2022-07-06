// Package cmd implements various PocketBase system commands.
package cmd

import (
	"fmt"

	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"
)

// NewVersionCommand creates and returns new command that prints
// the current PocketBase version.
func NewVersionCommand(app core.App, version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Prints the current PocketBase app version",
		Run: func(command *cobra.Command, args []string) {
			fmt.Printf("PocketBase v%s\n", version)
		},
	}
}
