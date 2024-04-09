package cli

import (
	"fmt"
	"os"
)

type Editor struct {
}

func NewEditor() *Editor {
	return &Editor{}
}

func (e *Editor) WriteToFile(filePath string, content string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	fmt.Printf("Content successfully written to: %s\n", filePath)

	return nil
}
