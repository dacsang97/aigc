package commit

import (
	"github.com/dacsang97/aigc/internal/provider"
)

type Generator struct {
	provider interface {
		Generate(changes, userMessage string, rules []string) (string, error)
	}
}

type ProviderConfig struct {
	Provider string
	Model    string
	APIKey   string
	Endpoint string
}

func New(config ProviderConfig) (*Generator, error) {
	p, err := provider.NewProvider(provider.Config{
		Provider: config.Provider,
		Model:    config.Model,
		APIKey:   config.APIKey,
		Endpoint: config.Endpoint,
	})
	if err != nil {
		return nil, err
	}

	return &Generator{
		provider: p,
	}, nil
}

func (g *Generator) Generate(changes, userMessage string, rules []string) (string, error) {
	return g.provider.Generate(changes, userMessage, rules)
}
