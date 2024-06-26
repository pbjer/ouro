package cli

import (
	"errors"
	"strings"

	"github.com/pbjer/ouro/internal/llm"
	"github.com/pbjer/ouro/internal/prompt"

	"gorm.io/gorm"

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

		_, err = NewPlanner(db).Plan(description)
		if err != nil {
			return err
		}

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

type Planner struct {
	db *gorm.DB
}

func NewPlanner(db *gorm.DB) *Planner {
	return &Planner{db}
}

func (p *Planner) GetLLMPreference() (string, error) {
	var config LLMConfig
	result := p.db.Last(&config) // Assuming you're storing the latest preference as the most recent entry.
	if result.Error != nil {
		return "", result.Error
	}
	return config.Provider, nil
}

func (p *Planner) Plan(description string) (*Plan, error) {
	var context []Context
	p.db.Find(&context)

	if len(context) == 0 {
		return nil, errors.New("failed to load context for plan")
	}

	var psString string
	err := projectStructure("./", "", &psString)
	if err != nil {
		return nil, err
	}
	planStartPrompt, err := prompt.New(prompt.PlannerCreatePlanPrompt, prompt.PlannerCreatePlanData{
		Change:           description,
		ProjectStructure: psString,
		Context:          NewContextPrompter(context...).Prompt(),
	}).Render()
	if err != nil {
		return nil, err
	}

	thread := llm.NewThread(
		llm.SystemMessage(prompt.PlannerSystemPrompt),
		llm.UserMessage(planStartPrompt),
	)

	clientChoice, err := p.GetLLMPreference()
	if err != nil {
		return nil, err
	}
	var clientOption llm.ClientOption
	switch clientChoice {
	case "openai":
		clientOption = llm.OpenAI
	case "groq":
		clientOption = llm.Groq
	}
	client := llm.NewClient(clientOption)

	err = client.Generate(thread)
	if err != nil {
		return nil, err
	}
	plan := Plan{
		Content: thread.LastMessage().Content,
	}
	p.db.Save(&plan)

	result := EditorFunction{}
	err = client.Map(plan.Content, &result)
	if err != nil {
		return nil, err
	}

	err = NewEditor().Edit(result.Edits...)
	if err != nil {
		return nil, err
	}

	return &plan, nil
}

type EditorFunction struct {
	Edits []Edit `json:"edits"`
}

type Edit struct {
	FilenameToChangeOrCreate string `json:"file_name_to_change_or_create"`
	CompleteFileContents     string `json:"complete_file_contents"`
}

type Plan struct {
	gorm.Model
	Content string
}
