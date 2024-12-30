package commit

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/dacsang97/aigc/cmd"
	"github.com/dacsang97/aigc/internal/commit"
	"github.com/dacsang97/aigc/internal/config"
	"github.com/dacsang97/aigc/internal/git"
	"github.com/dacsang97/aigc/internal/logger"
)

type Command struct {
	*cmd.BaseCommand
	configManager *config.Manager
	logger        *logger.Logger
	push          bool
}

func New(configManager *config.Manager, logger *logger.Logger) cmd.Command {
	c := &Command{
		configManager: configManager,
		logger:        logger,
	}

	baseCmd := cmd.NewBaseCommand(
		"commit",
		"Generate and create a commit",
		[]cli.Flag{
			&cli.BoolFlag{
				Name:        "push",
				Aliases:     []string{"p"},
				Usage:       "push changes after commit",
				Destination: &c.push,
			},
			&cli.StringFlag{
				Name:    "message",
				Aliases: []string{"m"},
				Usage:   "provide commit message hint (in any language)",
			},
		},
		c.handle,
	)

	c.BaseCommand = baseCmd
	return c
}

func (c *Command) handle(ctx *cli.Context) error {
	// Load local rules if they exist
	if err := c.configManager.LoadLocalRules(); err != nil {
		c.logger.DebugLog("Error loading local rules", err.Error())
	}

	// Initialize git client
	gitClient := git.New(c.push)

	// Get git changes
	changes, err := gitClient.GetStagedChanges()
	if err != nil {
		return err
	}

	c.logger.DebugLog("Git changes detected", changes)

	// Get user's commit message hint if provided
	userMessage := ctx.String("message")
	if userMessage != "" {
		c.logger.DebugLog("User provided commit message hint", userMessage)
	}

	// Initialize commit message generator
	generator := commit.New(c.configManager.Config.APIKey, c.configManager.Config.Model)

	// Generate commit message
	commitMsg, err := generator.Generate(changes, userMessage, c.configManager.GetRules())
	if err != nil {
		return err
	}

	c.logger.DebugLog("Generated commit message", commitMsg)

	// Commit changes
	if err := gitClient.Commit(commitMsg); err != nil {
		return err
	}

	fmt.Println("Successfully committed changes with message:")
	fmt.Println(commitMsg)

	return nil
}
