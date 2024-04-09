package cli

import (
	"fmt"
	"os"
	"path/filepath"
)

type Editor struct {
}

func NewEditor() *Editor {
	return &Editor{}
}

func (e *Editor) WriteToFile(filePath string, content string) error {
	// Extract the directory part from the filePath.
	dir := filepath.Dir(filePath)

	// Create all directories in the path (if they don't already exist).
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	// Proceed with file creation now that the path is ensured.
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	fmt.Println("Updated", filePath)

	return nil
}
