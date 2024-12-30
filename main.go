package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v2"
)

type Config struct {
	APIKey string `yaml:"api_key"`
	Model  string `yaml:"model"`
	Debug  bool   `yaml:"debug"`
}

var (
	config     Config
	configDir  string
	configPath string
	logDir     string
	debug      bool
	push       bool
	model      string
	logger     *zap.Logger
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	// Update config dir
	configDir = filepath.Join(home, ".aigc")
	configPath = filepath.Join(configDir, "config.yaml")
	// Create log dir
	logDir = filepath.Join(configDir, "log")
}

func loadConfig() error {
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("error creating config directory: %v", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			config = Config{
				Model: "google/gemini-flash-1.5-8b",
			}
			return nil
		}
		return err
	}

	return yaml.Unmarshal(data, &config)
}

func saveConfig() error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600)
}

func initLogger() {
	// Create log directory if not exists
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("failed to create log directory: %v", err)
	}

	// Get current date for log file name
	currentDate := time.Now().Format("2006-01-02")
	logFile := filepath.Join(logDir, currentDate+".log")

	// Create or append to log file
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}

	// Configure encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// Create core
	var core zapcore.Core
	if debug {
		// Debug mode: Log to both console and file
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

		consoleCore := zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			zap.DebugLevel,
		)
		fileCore := zapcore.NewCore(
			fileEncoder,
			zapcore.AddSync(file),
			zap.DebugLevel,
		)
		core = zapcore.NewTee(consoleCore, fileCore)
	} else {
		// Production mode: Log to file only
		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
		core = zapcore.NewCore(
			fileEncoder,
			zapcore.AddSync(file),
			zap.InfoLevel,
		)
	}

	// Create logger
	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
}

func debugLog(message string, content ...string) {
	if debug || config.Debug {
		if len(content) > 0 {
			logger.Debug(message,
				zap.String("content", content[0]),
			)
		} else {
			logger.Debug(message)
		}
	}
}

func main() {
	// Initialize logger first
	initLogger()
	defer logger.Sync()

	if err := loadConfig(); err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
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
		Commands: []*cli.Command{
			{
				Name:  "config",
				Usage: "Configure aicm settings",
				Flags: []cli.Flag{
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
				},
				Action: func(c *cli.Context) error {
					// Handle API key update
					if apiKey := c.String("apikey"); apiKey != "" {
						config.APIKey = apiKey
					}

					// Handle model update
					if model := c.String("model"); model != "" {
						config.Model = model
					}

					// Handle debug update
					if c.IsSet("debug") {
						config.Debug = c.Bool("debug")
					}

					// Save if any changes were made
					if c.NumFlags() > 0 {
						if err := saveConfig(); err != nil {
							return fmt.Errorf("error saving config: %v", err)
						}
						fmt.Println("Configuration updated successfully")
						return nil
					}

					// Show current config if no flags provided
					fmt.Printf("Current configuration:\n")
					fmt.Printf("API Key: %s\n", config.APIKey)
					fmt.Printf("Model: %s\n", config.Model)
					fmt.Printf("Debug: %v\n", config.Debug)
					return nil
				},
			},
			{
				Name:  "commit",
				Usage: "Generate and create a commit",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "push",
						Aliases:     []string{"p"},
						Usage:       "push changes after commit",
						Destination: &push,
					},
				},
				Action: runCommit,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal("application error", zap.Error(err))
	}
}
