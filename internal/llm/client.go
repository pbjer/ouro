package llm

import (
	"context"
	"fmt"
	"ouro/internal/env"

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
			Model:       openai.GPT4Turbo1106,
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
