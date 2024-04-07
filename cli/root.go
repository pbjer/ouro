package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "ouro",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Ouroboros started.")
	},
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
