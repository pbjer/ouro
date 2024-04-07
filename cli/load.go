package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
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
		db, err := InitDB() // Assume InitDB is implemented elsewhere
		if err != nil {
			return err
		}

		for _, path := range args {
			err := processPath(db, path)
			if err != nil {
				return fmt.Errorf("error processing path %s: %w", path, err)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)
	loadCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Recursively load files from directories")
	loadCmd.Flags().BoolVarP(&hidden, "hidden", "d", false, "Include hidden directories")

}

func processPath(db *gorm.DB, path string) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Printf("Path does not exist: %s\n", path)
		return nil // Continue processing other paths
	} else if err != nil {
		return err
	}

	if info.IsDir() {
		return processDirectory(db, path)
	} else {
		return processFile(db, path)
	}
}

func processDirectory(db *gorm.DB, dirPath string) error {
	return filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip any .ouro directory
		if strings.Contains(filepath.ToSlash(path), ".ouro") {
			return filepath.SkipDir
		}

		if strings.Contains(filepath.ToSlash(path), "/.") && !hidden {
			return filepath.SkipDir
		}

		if d.IsDir() {
			// If the directory is not the top-level directory being processed and -r is not set, skip it
			if path != dirPath && !recursive {
				return filepath.SkipDir
			}

			return nil
		}

		// Process each file
		return processFile(db, path)
	})
}

func processFile(db *gorm.DB, filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	var fileContent FileContent
	// Attempt to find an existing record and update it, or create a new one
	res := db.Where("path = ?", filePath).FirstOrCreate(&fileContent, FileContent{Path: filePath, Content: string(content)})
	if res.Error != nil {
		return res.Error
	}
	// If the file already exists, update its content
	if res.RowsAffected > 0 {
		fileContent.Content = string(content)
		db.Save(&fileContent)
	}

	fmt.Println("Loaded", filePath)

	return nil
}
