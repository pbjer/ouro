package llm

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
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

func SystemMessage(content string) Message {
	return NewMessage(RoleSystem, content)
}

func UserMessage(content string) Message {
	return NewMessage(RoleUser, content)
}

func AssistantMessage(content string) Message {
	return NewMessage(RoleAssistant, content)
}

func ToolMessage(content string) Message {
	return NewMessage(RoleTool, content)
}
