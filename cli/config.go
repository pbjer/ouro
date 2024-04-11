package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var configCmd = &cobra.Command{
	Use:   "use [provider]",
	Short: "Configure LLM preference",
	Args:  cobra.ExactArgs(1), // Ensures exactly one argument is passed
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := InitDB()
		if err != nil {
			return err
		}
		llmPreference := args[0] // Use the first argument as the LLM preference
		err = setLLMPreference(db, llmPreference)
		if err != nil {
			return err
		}

		fmt.Println("Using", llmPreference)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func setLLMPreference(db *gorm.DB, llmOption string) error {
	// First, check if the provided llmOption is valid
	llmOptions := []string{"openai", "groq"}
	isValidOption := false
	for _, option := range llmOptions {
		if llmOption == option {
			isValidOption = true
			break
		}
	}

	if !isValidOption {
		return fmt.Errorf("unsupported LLM option: %s", llmOption)
	}

	// Assuming the database only needs to track one global LLM setting
	config := LLMConfig{Provider: llmOption}
	result := db.Where("id = ?", 1).FirstOrCreate(&config) // Assuming you're managing a single config record
	if result.Error != nil {
		return result.Error
	}

	// If the record exists but the provider is different, update it
	if config.Provider != llmOption {
		config.Provider = llmOption
		if err := db.Save(&config).Error; err != nil {
			return err
		}
	}

	return nil
}

type LLMConfig struct {
	gorm.Model
	Provider string
}
