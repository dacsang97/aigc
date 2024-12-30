package provider

import "fmt"

// Provider represents an AI completion provider interface
type Provider interface {
	Generate(changes, userMessage string, rules []string) (string, error)
}

// Config represents the configuration for an AI provider
type Config struct {
	Provider string `yaml:"provider"` // "openai", "anthropic", "openrouter", or "custom"
	Model    string `yaml:"model"`    // The model to use
	APIKey   string `yaml:"api_key"`  // The API key for the provider
	Endpoint string `yaml:"endpoint"` // Custom API endpoint URL (optional)
}

// ProviderConfig contains provider-specific configurations
type ProviderConfig struct {
	BaseURL string
}

var (
	// SupportedProviders is a map of supported provider names to their configurations
	SupportedProviders = map[string]ProviderConfig{
		"openai": {
			BaseURL: "https://api.openai.com/v1/chat/completions",
		},
		"anthropic": {
			BaseURL: "https://api.anthropic.com/v1/messages",
		},
		"openrouter": {
			BaseURL: "https://openrouter.ai/api/v1/chat/completions",
		},
	}
)

// NewProvider creates a new AI provider based on the configuration
func NewProvider(config Config) (Provider, error) {
	if config.Provider == "custom" {
		if config.Endpoint == "" {
			return nil, fmt.Errorf("endpoint URL is required for custom provider")
		}
		return NewOpenAIProvider(Config{
			Provider: "custom",
			Model:    config.Model,
			APIKey:   config.APIKey,
			Endpoint: config.Endpoint,
		})
	}

	if config.Endpoint != "" {
		// If endpoint is provided, override the default BaseURL
		SupportedProviders[config.Provider] = ProviderConfig{BaseURL: config.Endpoint}
	}

	switch config.Provider {
	case "openai":
		return NewOpenAIProvider(config)
	case "anthropic":
		return NewAnthropicProvider(config)
	case "openrouter":
		return NewOpenRouterProvider(config)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", config.Provider)
	}
}
