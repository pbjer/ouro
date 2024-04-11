// reload.go
package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(reloadCmd)
}

var reloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Reload the to the latest version of what's currently loaded",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := InitDB()
		if err != nil {
			return fmt.Errorf("failed to initialize the database: %w", err)
		}

		files, err := ListFiles(db)
		if err != nil {
			return fmt.Errorf("failed to list file paths: %w", err)
		}

		// You will need to implement the LoadPath function in the context_loader.go
		for _, file := range files {
			loader := NewContextLoader(db)
			err := loader.LoadPath(file.Path)
			if err != nil {
				return fmt.Errorf("failed to reload path %s: %w", file.Path, err)
			}
		}

		paths, err := ListFiles(db)
		if err != nil {
			return err
		}

		totalTokens := 0
		for _, path := range paths {
			fmt.Println(path.Path, "-", path.Tokens)
			totalTokens += path.Tokens
		}

		fmt.Println("LOADED:", totalTokens)

		return nil
	},
}
