package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	projectStructureCmd = &cobra.Command{
		Use:   "ps",
		Short: "Provides information about the project structure",
		Run:   projectStructureRunner,
	}
	writeFlag  bool   // Flag to check if output should be written to a file
	outputFile string // The optional filename provided by the user
)

func init() {
	projectStructureCmd.Flags().BoolVarP(&writeFlag, "write", "w", false, "Write the output to a file instead of stdout. Uses PROJECT_STRUCTURE.md if no file is specified with --file")
	projectStructureCmd.Flags().StringVar(&outputFile, "file", "PROJECT_STRUCTURE.md", "Specifies the file name to write to when -w is used")
	rootCmd.AddCommand(projectStructureCmd)
}

func projectStructureRunner(cmd *cobra.Command, args []string) {
	directory := "./" // Default to current directory
	if len(args) == 1 {
		directory = args[0]
	}

	var result string
	err := projectStructure(directory, "", &result)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	if writeFlag {
		// Write to the specified file or default to PROJECT_STRUCTURE.md if not specified
		result = fmt.Sprintf("```\n" + result + "```\n")
		err = os.WriteFile(outputFile, []byte(result), 0644)
		if err != nil {
			fmt.Println("error writing to file:", err)
			return
		}
		fmt.Println("Written to", outputFile)
	} else {
		// Just print the result to stdout
		fmt.Print(result)
	}
}

var skipDirs = map[string]bool{
	"node_modules":    true,
	".git":            true,
	".DS_Store":       true,
	"vendor":          true,
	".idea":           true,
	"build":           true,
	"dist":            true,
	"out":             true,
	"tmp":             true,
	"temp":            true,
	".Trash":          true,
	".Spotlight-V100": true,
	".Trashes":        true,
}

func projectStructure(directory string, prefix string, result *string) error {
	files, err := os.ReadDir(directory)
	if err != nil {
		return err
	}

	for i, file := range files {
		if file.Name()[0] == '.' || skipDirs[file.Name()] {
			continue
		}

		newPrefix := prefix
		if i == len(files)-1 {
			*result += fmt.Sprintf("%s└── %s\n", prefix, file.Name())
			newPrefix += "    "
		} else {
			*result += fmt.Sprintf("%s├── %s\n", prefix, file.Name())
			newPrefix += "│   "
		}

		if file.IsDir() {
			err := projectStructure(filepath.Join(directory, file.Name()), newPrefix, result)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
