package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/gorm"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(unloadCmd)
}

var unloadCmd = &cobra.Command{
	Use:   "unload [paths...]",
	Short: "Unload stored file contents",
	Long:  `Removes records of file contents from the database for the specified files and directories, or all records if no paths are specified.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := InitDB() // Assume InitDB is implemented elsewhere
		if err != nil {
			return fmt.Errorf("failed to initialize the database: %w", err)
		}

		// Updated to pass args to unloadFileContents
		if err := unloadFileContents(db, args...); err != nil {
			return fmt.Errorf("failed to unload file contents: %w", err)
		}

		return nil
	},
}

func unloadFileContents(db *gorm.DB, paths ...string) error {
	if len(paths) > 0 {
		for _, path := range paths {
			// Determine if the path is a directory or a file
			fileInfo, err := os.Stat(path)
			if err != nil {
				fmt.Printf("Unable to access %s: %s\n", path, err)
				continue // Skip this path if it cannot be accessed
			}

			if fileInfo.IsDir() {
				// If it's a directory, fetch entries for all files within this directory before deleting
				cleanedPath := filepath.Clean(path) + "/" // Ensure matching with subpaths
				var filesToDelete []Context
				if err := db.Where("path LIKE ?", fmt.Sprintf("%s%%", cleanedPath)).Find(&filesToDelete).Error; err != nil {
					return err
				}
				if err := db.Where("path LIKE ?", fmt.Sprintf("%s%%", cleanedPath)).Delete(&Context{}).Error; err != nil {
					return err
				}
				for _, file := range filesToDelete {
					fmt.Println("Unloaded", file.Path)
				}
			} else {
				// It's a file, delete the specific entry
				if err := db.Where("path = ?", path).Delete(&Context{}).Error; err != nil {
					return err
				}
				fmt.Println("Unloaded", path)
			}
		}
	} else {
		// Retrieve all records before deleting them
		var allFileContents []Context
		if err := db.Find(&allFileContents).Error; err != nil {
			return fmt.Errorf("failed to retrieve file contents: %w", err)
		}

		// Delete all records
		if err := db.Where("1=1").Delete(&Context{}).Error; err != nil {
			return fmt.Errorf("failed to unload file contents: %w", err)
		}

		// Print paths of all unloaded contents
		if len(allFileContents) > 0 {
			for _, fileContent := range allFileContents {
				fmt.Println("Unloaded", fileContent.Path)
			}
		} else {
			fmt.Println("Nothing to unload.")
		}
	}
	return nil
}
