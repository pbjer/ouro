package llm

import "strings"

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type Message struct {
	Role    Role
	Content string
}

func NewMessage(role Role, content string) Message {
	return Message{
		Role:    role,
		Content: content,
	}
}

func (m *Message) String() string {
	fmtRole := strings.ToUpper(string(m.Role))
	return "[" + fmtRole + "]" + "-----------------------------------\n" + m.Content
}

func SystemMessage(content string) Message {
	return NewMessage(RoleSystem, content)
}

func UserMessage(content string) Message {
	return NewMessage(RoleUser, content)
}

func AssistantMessage(content string) Message {
	return NewMessage(RoleAssistant, content)
}
