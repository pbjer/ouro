package llm

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
