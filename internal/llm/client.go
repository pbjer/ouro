package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/pbjer/ouro/internal/env"

	"github.com/sashabaranov/go-openai"
)

type Client struct {
	client *openai.Client
}

func NewClient() *Client {
	config := openai.DefaultConfig(env.GroqAPIKey())
	config.BaseURL = "https://api.groq.com/openai/v1"
	client := openai.NewClientWithConfig(config)
	return &Client{
		client: client,
	}
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

	prompt := fmt.Sprintf("Given the source data:\n\n%s\n\nProvide a JSON object filled out with the values from the source data that correspond to these fields: %s.", source, fields)
	prompt += "\nIMPORTANT: DO NOT WRAP THE RESPONSE IN MARKDOWN, IT MUST BE RAW JSON!"
	prompt += "\nIMPORTANT: THE FIELDS MUST APPEAR IN THE JSON EXACTLY AS WRITTEN ABOVE!"
	prompt += "\nIMPORTANT: YOU MUST PROPERLY ESCAPE ALL QUOTES WITHIN THE CONTENT SO THAT THE JSON IS VALID!"

	thread := NewThread(SystemMessage(prompt))
	fmt.Println(thread.String())
	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT4Turbo1106,
			Messages:    ThreadToOpenAICompletionMessages(thread),
			Temperature: 0,
		},
	)
	if err != nil {
		return err
	}

	lastMessage := resp.Choices[0].Message.Content

	// Ensure we only have JSON format in the message.
	if !json.Valid([]byte(lastMessage)) {
		cleaned := TrimNonJSON(lastMessage)
		if !json.Valid([]byte(cleaned)) {
			return fmt.Errorf("json still invalid after cleaning")
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
