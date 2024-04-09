package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"ouro/internal/env"
	"reflect"

	"github.com/sashabaranov/go-openai"
)

type Client struct {
	client *openai.Client
}

func NewClient() *Client {
	return &Client{
		client: openai.NewClient(env.OpenAIAPIKey()),
	}
}

func (c *Client) Generate(thread *Thread) error {
	numTokens, err := thread.NumberOfTokens()
	if err != nil {
		return err
	}
	fmt.Printf("\n\nGeneration request tokens: %d\n\n", numTokens)
	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo,
			Messages:    ThreadToOpenAICompletionMessages(thread),
			Temperature: 0,
		},
	)
	if err != nil {
		return err
	}
	thread.AddMessages(AssistantMessage(resp.Choices[0].Message.Content))
	return nil
}

func (c *Client) Map(source string, target interface{}) error {
	targetType := reflect.TypeOf(target)
	if targetType.Kind() != reflect.Ptr || targetType.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to a struct")
	}
	structType := targetType.Elem()
	fields := make([]string, 0, structType.NumField())
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		jsonTag, ok := field.Tag.Lookup("json")
		if !ok {
			continue
		}
		fields = append(fields, jsonTag)
	}

	prompt := fmt.Sprintf("Given the source data:\n\n%s\n\nProvide a JSON object with the following fields: %s.", source, fields)
	prompt += "\nIMPORTANT: DO NOT WRAP THE RESPONSE IN MARKDOWN, IT MUST BE RAW JSON!"
	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4Turbo1106,
			Messages: []openai.ChatCompletionMessage{
				{Role: string(RoleSystem), Content: prompt},
			},
			Temperature: 0,
		},
	)
	if err != nil {
		return err
	}

	lastMessage := resp.Choices[0].Message.Content

	// Ensure we only have JSON format in the message.
	if !json.Valid([]byte(lastMessage)) {
		fmt.Println(lastMessage)
		cleaned := TrimNonJSON(lastMessage)
		if !json.Valid([]byte(cleaned)) {
			return fmt.Errorf("")
		}
		return fmt.Errorf("response from LLM is not valid JSON")
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
