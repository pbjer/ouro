package llm

import "strings"

type Thread struct {
	Messages []Message
}

func NewThread(messages ...Message) *Thread {
	return &Thread{
		Messages: messages,
	}
}

func (t *Thread) AddMessages(messages ...Message) {
	t.Messages = append(t.Messages, messages...)
}

func (t *Thread) LastMessage() Message {
	numberOfMessages := len(t.Messages)
	if numberOfMessages == 0 {
		return Message{}
	}
	return t.Messages[numberOfMessages-1]
}

// NumberOfTokens computes the total number of tokens in all messages of the thread
func (t *Thread) NumberOfTokens() (int, error) {
	totalTokens := 0
	for _, message := range t.Messages {
		tokens, err := NumberOfTokens(message.Content)
		if err != nil {
			return 0, err
		}
		totalTokens += tokens
	}
	return totalTokens, nil
}

func (t *Thread) String() string {
	var sb strings.Builder
	for _, message := range t.Messages {
		sb.WriteString("\n" + message.String())
	}
	sb.WriteString("\n")
	return sb.String()
}
