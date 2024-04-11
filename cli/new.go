package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Initialize a new ouro project",
	Long:  `This command initializes a new ouro project in the current directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := InitDB()
		if err != nil {
			return err
		}
		fmt.Println("Initialized ouro")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func InitDB() (*gorm.DB, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	ouroDir := filepath.Join(workingDir, ".ouro")
	if _, err := os.Stat(ouroDir); os.IsNotExist(err) {
		if err := os.Mkdir(ouroDir, os.ModePerm); err != nil {
			return nil, err
		}
	}

	dbPath := filepath.Join(ouroDir, "sqlite.db")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	err = db.AutoMigrate(&Context{}, &Plan{}, &LLMConfig{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
