package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dacsang97/aigc/internal/prompt"
)

type OpenRouterProvider struct {
	config  Config
	baseURL string
	prompt  *prompt.Generator
}

func NewOpenRouterProvider(config Config) (*OpenRouterProvider, error) {
	baseURL := config.Endpoint
	if baseURL == "" {
		baseURL = SupportedProviders["openrouter"].BaseURL
	}

	return &OpenRouterProvider{
		config:  config,
		baseURL: baseURL,
		prompt:  prompt.New(),
	}, nil
}

type OpenRouterRequestBody struct {
	Stream   bool             `json:"stream"`
	Model    string           `json:"model"`
	Messages []prompt.Message `json:"messages"`
}

func (p *OpenRouterProvider) Generate(changes, userMessage string, rules []string) (string, error) {
	messages := p.prompt.BuildMessages(changes, userMessage, rules)

	reqBody := OpenRouterRequestBody{
		Stream:   false,
		Model:    p.config.Model,
		Messages: messages,
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

	req.Header.Set("Authorization", "Bearer "+p.config.APIKey)
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

	return apiResp.Choices[0].Message.Content, nil
}
