package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"

	"github.com/urfave/cli/v2"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestBody struct {
	Stream   bool      `json:"stream"`
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Choice struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
}

type APIResponse struct {
	Choices []Choice `json:"choices"`
}

func generateCommitMessage(changes string, userMessage string) (string, error) {
	// Load local rules if they exist
	if err := loadLocalRules(); err != nil {
		debugLog("Error loading local rules", err.Error())
	}

	// Build system message with rules
	systemContent := "Generate a commit message following the Conventional Commits standard:\n\n" +
		"<type>[optional scope]: <description>\n\n" +
		"[optional body]\n\n" +
		"[optional footer]\n\n" +
		"Rules:\n" +
		"1. Type MUST be one of:\n" +
		"   - feat: new feature (correlates with MINOR version)\n" +
		"   - fix: bug fix (correlates with PATCH version)\n" +
		"   - docs: documentation changes\n" +
		"   - style: formatting, missing semi colons, etc\n" +
		"   - refactor: refactoring code\n" +
		"   - perf: performance improvements\n" +
		"   - test: adding tests\n" +
		"   - chore: maintenance tasks\n\n" +
		"2. Scope is optional and should describe the section of code (e.g., feat(parser))\n" +
		"3. Description must be concise and in imperative mood (e.g., 'change' not 'changed')\n" +
		"4. Body should explain the motivation for the change and contrast with previous behavior\n" +
		"5. Breaking changes MUST be indicated by BREAKING CHANGE: in footer\n" +
		"6. A ! MAY be added before the : for breaking changes (e.g., feat!: breaking change)\n\n"

	// Add project-specific rules if any
	if len(config.Rules) > 0 {
		systemContent += "Project-specific rules:\n"
		for _, rule := range config.Rules {
			systemContent += fmt.Sprintf("- %s\n", rule)
		}
		systemContent += "\n"
	}

	systemContent += "IMPORTANT: Always generate the commit message in English, regardless of the input language.\n" +
		"Do not include any explanation in your response, only return the commit message content."

	messages := []Message{
		{
			Role:    "system",
			Content: systemContent,
		},
	}

	// Add user's commit message hint if provided
	if userMessage != "" {
		messages = append(messages, Message{
			Role: "user",
			Content: fmt.Sprintf("User provided this commit message hint (which may be in any language):\n%s\n\n"+
				"Please consider this message when generating the commit message. "+
				"Understand the meaning and translate the intent to English if needed, "+
				"but ensure the output follows the Conventional Commits standard and is in English.",
				userMessage),
		})
	}

	// Add the code changes
	messages = append(messages, Message{
		Role: "user",
		Content: fmt.Sprintf("Analyze these file changes and generate a commit message:\n```\n%s\n```\n\n"+
			"Guidelines:\n"+
			"1. Use appropriate type based on the changes (feat for new features, fix for bugs, etc.)\n"+
			"2. Add relevant scope if the changes are focused on a specific component\n"+
			"3. Description must be under 100 characters\n"+
			"4. Include breaking changes in footer with BREAKING CHANGE: prefix if any\n"+
			"5. Add detailed body explaining motivation and changes if significant\n"+
			"6. Use issue/PR references in footer if relevant\n\n"+
			"Return only the commit message without any extra content or backticks.",
			changes),
	})

	reqBody := RequestBody{
		Stream:   false,
		Model:    config.Model,
		Messages: messages,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	if config.APIKey == "" {
		return "", fmt.Errorf("API key not found. Please run 'aicm config <your_api_key>' first")
	}

	req.Header.Set("Authorization", "Bearer "+config.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", err
	}

	if len(apiResp.Choices) == 0 {
		return "", fmt.Errorf("no commit message generated")
	}

	// Clean and return the commit message
	return apiResp.Choices[0].Message.Content, nil
}

func runCommit(c *cli.Context) error {
	// Update config if flags are set
	if c.String("model") != "" {
		config.Model = c.String("model")
	}
	config.Debug = c.Bool("debug")

	// Get git changes
	changes, err := getGitChanges()
	if err != nil {
		return err
	}

	debugLog("Git changes detected", changes)

	// Get user's commit message hint if provided
	userMessage := c.String("message")
	if userMessage != "" {
		debugLog("User provided commit message hint", userMessage)
	}

	// Generate commit message
	commitMsg, err := generateCommitMessage(changes, userMessage)
	if err != nil {
		return err
	}

	debugLog("Generated commit message", commitMsg)

	// Commit changes
	if err := commitChanges(commitMsg); err != nil {
		return err
	}

	fmt.Println("Successfully committed changes with message:")
	fmt.Println(commitMsg)

	// Push if requested
	if push {
		if err := pushChanges(); err != nil {
			return err
		}
		fmt.Println("Successfully pushed changes")
	}

	return nil
}

func pushChanges() error {
	cmd := exec.Command("git", "push")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to push changes: %v\n%s", err, output)
	}
	return nil
}
