package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dacsang97/aigc/internal/prompt"
)

type AnthropicProvider struct {
	config  Config
	baseURL string
	prompt  *prompt.Generator
}

func NewAnthropicProvider(config Config) (*AnthropicProvider, error) {
	baseURL := config.Endpoint
	if baseURL == "" {
		baseURL = SupportedProviders["anthropic"].BaseURL
	}

	return &AnthropicProvider{
		config:  config,
		baseURL: baseURL,
		prompt:  prompt.New(),
	}, nil
}

type AnthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AnthropicRequestBody struct {
	Model     string             `json:"model"`
	Messages  []AnthropicMessage `json:"messages"`
	MaxTokens int                `json:"max_tokens"`
}

type AnthropicResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

func (p *AnthropicProvider) Generate(changes, userMessage string, rules []string) (string, error) {
	messages := p.prompt.BuildMessages(changes, userMessage, rules)
	anthropicMessages := make([]AnthropicMessage, len(messages))
	for i, msg := range messages {
		anthropicMessages[i] = AnthropicMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	reqBody := AnthropicRequestBody{
		Model:     p.config.Model,
		Messages:  anthropicMessages,
		MaxTokens: 1000,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", p.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	if p.config.APIKey == "" {
		return "", fmt.Errorf("API key not found. Please run 'aicm config' to set up your configuration")
	}

	req.Header.Set("x-api-key", p.config.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")
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

	var apiResp AnthropicResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", err
	}

	if len(apiResp.Content) == 0 {
		return "", fmt.Errorf("no commit message generated")
	}

	return apiResp.Content[0].Text, nil
}
