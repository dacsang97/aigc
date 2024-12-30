package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/dacsang97/aigc/cmd"
	cmdcommit "github.com/dacsang97/aigc/cmd/commit"
	cmdconfig "github.com/dacsang97/aigc/cmd/config"
	"github.com/dacsang97/aigc/internal/config"
	"github.com/dacsang97/aigc/internal/logger"
)

var (
	configManager *config.Manager
	appLogger     *logger.Logger
	debug         bool
	model         string
)

func init() {
	var err error
	configManager, err = config.NewManager()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if err := configManager.Load(); err != nil {
		log.Fatal("failed to load config:", err)
	}

	var err error
	appLogger, err = logger.New(configManager.LogDir, debug || configManager.Config.Debug)
	if err != nil {
		log.Fatal("failed to initialize logger:", err)
	}
	defer appLogger.Sync()

	// Initialize commands
	commands := []cmd.Command{
		cmdconfig.New(configManager, appLogger),
		cmdcommit.New(configManager, appLogger),
	}

	// Convert commands to cli.Commands
	cliCommands := make([]*cli.Command, len(commands))
	for i, command := range commands {
		cliCommands[i] = cmd.ToCLICommand(command)
	}

	app := &cli.App{
		Name:    "aigc",
		Usage:   "AI-powered Git commit message generator",
		Version: "1.0.0",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "debug",
				Usage:       "enable debug mode",
				Destination: &debug,
			},
			&cli.StringFlag{
				Name:        "model",
				Usage:       "AI model to use",
				Value:       "google/gemini-flash-1.5-8b",
				Destination: &model,
			},
		},
		Commands: cliCommands,
	}

	if err := app.Run(os.Args); err != nil {
		appLogger.Fatal("application error", zap.Error(err))
	}
}
