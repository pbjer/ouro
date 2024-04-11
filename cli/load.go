package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var recursive bool // Flag for recursive directory processing
var hidden bool

var loadCmd = &cobra.Command{
	Use:   "load [flags] [files and directories]",
	Short: "Load the contents of specified files and directories into the database",
	Long: `Load the contents of specified files and directories into the database. 
		If a directory is provided, files within are loaded non-recursively by default. 
		Use the -r flag for recursive loading.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := InitDB()
		if err != nil {
			return err
		}

		loader := NewContextLoader(db)
		for _, path := range args {
			err := loader.LoadPath(path)
			if err != nil {
				return fmt.Errorf("error processing path %s: %w", path, err)
			}
		}

		paths, err := ListFiles(db)
		if err != nil {
			return err
		}

		for _, path := range paths {
			fmt.Println("Loaded", path.Path, "-", path.Tokens)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)
	loadCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Recursively load files from directories")
	loadCmd.Flags().BoolVarP(&hidden, "hidden", "d", false, "Include hidden directories")
}
