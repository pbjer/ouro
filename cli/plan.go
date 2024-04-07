package cli

import (
	"errors"
	"fmt"
	"ouro/internal/llm"
	"ouro/internal/prompt"
	"strings"

	"github.com/spf13/cobra"
)

// planCmd represents the plan command
var planCmd = &cobra.Command{
	Use:   "plan [description]",
	Short: "Accepts a description of goals and creates a plan to accomplish them.",
	Args:  cobra.ExactArgs(1), // Ensures exactly one argument is passed
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := InitDB()
		if err != nil {
			return err
		}

		description := args[0]

		thread := llm.NewThread(
			llm.SystemMessage(prompt.PlannerSystemPrompt))

		var context []Context
		db.Find(&context)

		if len(context) == 0 {
			return errors.New("failed to load context for plan")
		}

		createPlan, err := prompt.New(prompt.PlannerCreatePlanPrompt, prompt.PlannerCreatePlanData{
			Change:  description,
			Context: NewContextPrompter(context...).Prompt(),
		}).Render()
		if err != nil {
			return err
		}

		thread.AddMessages(llm.UserMessage(createPlan))

		err = llm.NewClient().Generate(thread)
		if err != nil {
			return err
		}

		fmt.Println("Plan:")
		fmt.Println(thread.LastMessage().Content)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
}

type ContextPrompter struct {
	Context []Context
}

func (p *ContextPrompter) Prompt() string {
	var sb strings.Builder
	for _, context := range p.Context {
		sb.WriteString("\n" + context.Path + ":\n")
		sb.WriteString("```\n")
		sb.WriteString(context.Content)
		sb.WriteString("\n```")
	}
	return sb.String()
}

func NewContextPrompter(context ...Context) *ContextPrompter {
	return &ContextPrompter{context}
}
