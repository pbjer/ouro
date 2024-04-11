package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all loaded file paths",
	Long:  `Lists all the file paths that have been loaded into the database.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := InitDB()
		if err != nil {
			return err
		}

		paths, err := ListFiles(db)
		if err != nil {
			return err
		}

		if len(paths) == 0 {
			fmt.Println("Nothing to list.")
		}

		for _, path := range paths {
			fmt.Println(path.Path, "-", path.Tokens)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func ListFiles(db *gorm.DB) ([]Context, error) {
	var files []Context
	if result := db.Find(&files); result.Error != nil {
		return nil, result.Error
	}
	return files, nil
}
