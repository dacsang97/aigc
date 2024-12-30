package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Provider struct {
		Provider string `yaml:"provider"` // "openai" or "openrouter" or "custom"
		Model    string `yaml:"model"`    // The model to use
		APIKey   string `yaml:"api_key"`  // The API key for the provider
		Endpoint string `yaml:"endpoint"` // Custom API endpoint URL (optional)
	} `yaml:"provider"`
	Debug bool   `yaml:"debug"`
	Rules string `yaml:"rules"`
}

type Manager struct {
	Config     Config
	ConfigDir  string
	ConfigPath string
	LogDir     string
}

func NewManager() (*Manager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configDir := filepath.Join(home, ".aigc")
	configPath := filepath.Join(configDir, "config.yaml")
	logDir := filepath.Join(configDir, "log")

	return &Manager{
		ConfigDir:  configDir,
		ConfigPath: configPath,
		LogDir:     logDir,
	}, nil
}

func (m *Manager) Load() error {
	configDir := filepath.Dir(m.ConfigPath)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("error creating config directory: %v", err)
	}

	data, err := os.ReadFile(m.ConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			m.Config = Config{
				Provider: struct {
					Provider string `yaml:"provider"`
					Model    string `yaml:"model"`
					APIKey   string `yaml:"api_key"`
					Endpoint string `yaml:"endpoint"`
				}{
					Provider: "openrouter",
					Model:    "google/gemini-flash-1.5-8b",
					Endpoint: "",
				},
			}
			return nil
		}
		return err
	}

	return yaml.Unmarshal(data, &m.Config)
}

func (m *Manager) Save() error {
	data, err := yaml.Marshal(m.Config)
	if err != nil {
		return err
	}

	configDir := filepath.Dir(m.ConfigPath)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return err
	}

	return os.WriteFile(m.ConfigPath, data, 0600)
}

func (m *Manager) LoadLocalRules() error {
	data, err := os.ReadFile(".aigcrules")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if m.Config.Rules != "" {
		m.Config.Rules += "\n"
	}
	m.Config.Rules += strings.TrimSpace(string(data))
	return nil
}

func (m *Manager) GetRules() []string {
	if m.Config.Rules == "" {
		return nil
	}
	return strings.Split(strings.TrimSpace(m.Config.Rules), "\n")
}
