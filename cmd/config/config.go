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
		"Configure aicm settings",
		[]cli.Flag{
			&cli.StringFlag{
				Name:  "apikey",
				Usage: "Set API key",
			},
			&cli.StringFlag{
				Name:  "model",
				Usage: "Set AI model",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Set debug mode",
			},
			&cli.StringFlag{
				Name:  "rules",
				Usage: "Set commit message generation rules (one rule per line)",
			},
		},
		c.handle,
	)

	c.BaseCommand = baseCmd
	return c
}

func (c *Command) handle(ctx *cli.Context) error {
	// Handle API key update
	if apiKey := ctx.String("apikey"); apiKey != "" {
		c.configManager.Config.APIKey = apiKey
	}

	// Handle model update
	if model := ctx.String("model"); model != "" {
		c.configManager.Config.Model = model
	}

	// Handle debug update
	if ctx.IsSet("debug") {
		c.configManager.Config.Debug = ctx.Bool("debug")
	}

	// Handle rules update
	if rules := ctx.String("rules"); rules != "" {
		c.configManager.Config.Rules = rules
	}

	// Save if any changes were made
	if ctx.NumFlags() > 0 {
		if err := c.configManager.Save(); err != nil {
			return fmt.Errorf("error saving config: %v", err)
		}
		fmt.Println("Configuration updated successfully")
		return nil
	}

	// Show current config if no flags provided
	fmt.Printf("Current configuration:\n")
	fmt.Printf("API Key: %s\n", c.configManager.Config.APIKey)
	fmt.Printf("Model: %s\n", c.configManager.Config.Model)
	fmt.Printf("Debug: %v\n", c.configManager.Config.Debug)
	rules := c.configManager.GetRules()
	if len(rules) > 0 {
		fmt.Printf("\nCommit message rules:\n")
		for _, rule := range rules {
			fmt.Printf("  - %s\n", strings.TrimSpace(rule))
		}
	}
	return nil
}
