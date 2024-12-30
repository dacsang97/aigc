package config

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/dacsang97/aigc/cmd"
	"github.com/dacsang97/aigc/internal/config"
	"github.com/dacsang97/aigc/internal/logger"
)

type Command struct {
	*cmd.BaseCommand
	configManager *config.Manager
	logger        *logger.Logger
}

func New(configManager *config.Manager, logger *logger.Logger) cmd.Command {
	c := &Command{
		configManager: configManager,
		logger:        logger,
	}

	baseCmd := cmd.NewBaseCommand(
		"config",
		"Configure the application settings",
		[]cli.Flag{
			&cli.StringFlag{
				Name:  "provider",
				Usage: "Set the AI provider (openai, openrouter, or custom)",
			},
			&cli.StringFlag{
				Name:  "model",
				Usage: "Set the AI model",
			},
			&cli.StringFlag{
				Name:  "api-key",
				Usage: "Set the API key",
			},
			&cli.StringFlag{
				Name:  "endpoint",
				Usage: "Set custom API endpoint URL (optional)",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Enable debug mode",
			},
		},
		c.runConfig,
	)

	c.BaseCommand = baseCmd
	return c
}

func (c *Command) runConfig(ctx *cli.Context) error {
	if err := c.configManager.Load(); err != nil {
		return err
	}

	updated := false

	if provider := ctx.String("provider"); provider != "" {
		provider = strings.ToLower(provider)
		if provider != "openai" && provider != "openrouter" {
			return fmt.Errorf("invalid provider: %s. Must be 'openai' or 'openrouter'", provider)
		}
		c.configManager.Config.Provider.Provider = provider
		updated = true
	}

	if model := ctx.String("model"); model != "" {
		c.configManager.Config.Provider.Model = model
		updated = true
	}

	if apiKey := ctx.String("api-key"); apiKey != "" {
		c.configManager.Config.Provider.APIKey = apiKey
		updated = true
	}

	if endpoint := ctx.String("endpoint"); endpoint != "" {
		c.configManager.Config.Provider.Endpoint = endpoint
		updated = true
	}

	if ctx.IsSet("debug") {
		c.configManager.Config.Debug = ctx.Bool("debug")
		updated = true
	}

	if !updated {
		fmt.Printf("Current configuration:\n")
		fmt.Printf("  Provider: %s\n", c.configManager.Config.Provider.Provider)
		fmt.Printf("  Model: %s\n", c.configManager.Config.Provider.Model)
		fmt.Printf("  API Key: %s\n", maskAPIKey(c.configManager.Config.Provider.APIKey))
		fmt.Printf("  Endpoint: %s\n", c.configManager.Config.Provider.Endpoint)
		fmt.Printf("  Debug: %v\n", c.configManager.Config.Debug)
		return nil
	}

	if err := c.configManager.Save(); err != nil {
		return err
	}

	c.logger.Info("Configuration updated successfully")
	return nil
}

func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return "********"
	}
	return key[:4] + "..." + key[len(key)-4:]
}
