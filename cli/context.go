package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkoukk/tiktoken-go"
	"gorm.io/gorm"
)

type Context struct {
	gorm.Model
	Path    string
	Content string `gorm:"type:text;"`
	Tokens  int
}

type ContextLoader struct {
	db *gorm.DB
}

func NewContextLoader(db *gorm.DB) *ContextLoader {
	return &ContextLoader{
		db: db,
	}
}

func (l *ContextLoader) LoadPath(path string) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Printf("Path does not exist: %s\n", path)
		return nil // Continue processing other paths
	} else if err != nil {
		return err
	}

	if info.IsDir() {
		return l.processDirectory(path)
	} else {
		return l.processFile(path)
	}
}

func (l *ContextLoader) processDirectory(dirPath string) error {
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
		return l.processFile(path)
	})
}

func tokenCount(text string) (int, error) {
	tkm, err := tiktoken.EncodingForModel("gpt-4")
	if err != nil {
		err = fmt.Errorf("error getting encoding for model: %v", err)
		return 0, err
	}
	return len(tkm.Encode(text, nil, nil)), nil
}

func (l *ContextLoader) processFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	tokenCount, err := tokenCount(string(content))
	if err != nil {
		return err
	}

	memory := Context{
		Tokens:  tokenCount,
		Path:    filePath,
		Content: string(content),
	}

	// Attempt to find an existing record and update it, or create a new one
	res := l.db.Where("path = ?", filePath).FirstOrCreate(&memory, memory)
	if res.Error != nil {
		return res.Error
	}
	// If the file already exists, update its content
	if res.RowsAffected > 0 {
		l.db.Save(&memory)
	}

	fmt.Println("Loaded", filePath)

	return nil
}
