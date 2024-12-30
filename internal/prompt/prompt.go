package prompt

import "fmt"

// Generator handles the generation of prompts for AI models
type Generator struct {
	systemTemplate string
	userTemplate   string
}

// New creates a new prompt generator with default templates
func New() *Generator {
	return &Generator{
		systemTemplate: defaultSystemTemplate,
		userTemplate:   defaultUserTemplate,
	}
}

// BuildMessages builds the complete message list for the AI model
func (g *Generator) BuildMessages(changes, userMessage string, rules []string) []Message {
	messages := []Message{
		{
			Role:    "system",
			Content: g.buildSystemMessage(rules),
		},
	}

	if userMessage != "" {
		messages = append(messages, Message{
			Role:    "user",
			Content: g.buildUserHintMessage(userMessage),
		})
	}

	messages = append(messages, Message{
		Role:    "user",
		Content: g.buildChangeMessage(changes),
	})

	return messages
}

func (g *Generator) buildSystemMessage(rules []string) string {
	systemContent := g.systemTemplate

	if len(rules) > 0 {
		systemContent += "\nProject-specific rules:\n"
		for _, rule := range rules {
			systemContent += fmt.Sprintf("- %s\n", rule)
		}
		systemContent += "\n"
	}

	return systemContent
}

func (g *Generator) buildUserHintMessage(hint string) string {
	return fmt.Sprintf(g.userTemplate, hint)
}

func (g *Generator) buildChangeMessage(changes string) string {
	return fmt.Sprintf(defaultChangeMessageTemplate, changes)
}
