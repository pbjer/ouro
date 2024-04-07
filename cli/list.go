package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

// Define the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all loaded file paths",
	Long:  `Lists all the file paths that have been loaded into the database.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := InitDB()
		if err != nil {
			return err
		}

		paths, err := ListFilePaths(db)
		if err != nil {
			return err
		}

		for _, path := range paths {
			fmt.Println(path)
		}

		if len(paths) == 0 {
			fmt.Println("Nothing to list.")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

// New function to list all file paths
func ListFilePaths(db *gorm.DB) ([]string, error) {
	var files []FileContent
	if result := db.Find(&files); result.Error != nil {
		return nil, result.Error
	}

	paths := make([]string, len(files))
	for i, file := range files {
		paths[i] = file.Path
	}

	return paths, nil
}
