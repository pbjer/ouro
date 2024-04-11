package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/pbjer/ouro/internal/prompt"

	"github.com/pbjer/ouro/internal/env"

	"github.com/sashabaranov/go-openai"
)

type Client struct {
	client  *openai.Client
	baseURL string
	model   string
	apiKey  string
}

type ClientOption func(c *Client)

func Groq(client *Client) {
	client.baseURL = "https://api.groq.com/openai/v1"
	client.apiKey = env.GroqAPIKey()
}

func OpenAI(client *Client) {
	client.baseURL = "https://api.openai.com/v1/chat/completions"
	client.apiKey = env.OpenAIAPIKey()
}

func NewClient(options ...ClientOption) *Client {
	config := openai.DefaultConfig(env.GroqAPIKey())
	config.BaseURL = "https://api.groq.com/openai/v1"
	openaiclient := openai.NewClientWithConfig(config)
	client := &Client{
		client: openaiclient,
	}
	for _, option := range options {
		option(client)
	}
	return client
}

func (c *Client) Generate(thread *Thread) error {
	numTokens, err := thread.NumberOfTokens()
	if err != nil {
		return err
	}
	fmt.Println("Request tokens:", numTokens)
	fmt.Println(thread.String())
	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       "mixtral-8x7b-32768",
			Messages:    ThreadToOpenAICompletionMessages(thread),
			Temperature: 0,
		},
	)
	if err != nil {
		return err
	}
	result := AssistantMessage(resp.Choices[0].Message.Content)
	fmt.Println(result)
	thread.AddMessages(result)
	return nil
}

func (c *Client) Map(source string, target interface{}) error {
	targetType := reflect.TypeOf(target)
	if targetType.Kind() != reflect.Ptr || targetType.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to a struct")
	}

	jsonString, err := StructToJSON(target)
	if err != nil {
		return err
	}

	mappingPrompt, err := prompt.New(prompt.JSONMappingPromptTemplate, prompt.JSONMappingPromptData{
		Source:      source,
		ExampleJSON: jsonString,
	}).Render()
	if err != nil {
		return err
	}

	thread := NewThread(SystemMessage(mappingPrompt))
	err = c.Generate(thread)
	if err != nil {
		return err
	}

	lastMessage := thread.LastMessage().Content

	fmt.Println(lastMessage)

	// Ensure we only have JSON format in the message.
	if !json.Valid([]byte(lastMessage)) {
		trimmed := TrimNonJSON(lastMessage)
		cleaned := CleanJSONString(trimmed)
		if !json.Valid([]byte(cleaned)) {
			return json.Unmarshal([]byte(lastMessage), target)
		}
		lastMessage = cleaned
	}

	return json.Unmarshal([]byte(lastMessage), target)
}

func ThreadToOpenAICompletionMessages(thread *Thread) (messages []openai.ChatCompletionMessage) {
	for _, message := range thread.Messages {
		messages = append(messages, MessageToOpenAICompletionMessage(message))
	}
	return
}

func MessageToOpenAICompletionMessage(message Message) openai.ChatCompletionMessage {
	return openai.ChatCompletionMessage{
		Role:    string(message.Role),
		Content: message.Content,
	}
}

func CleanJSONString(s string) string {
	// Replace unescaped tabs with \\t (the escaped tab in JSON).
	cleaned := strings.Replace(s, "\t", "\\t", -1)

	return cleaned
}
