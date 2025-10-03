package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the configuration loaded from YAML
type Config struct {
	Credentials struct {
		OpenAIAPIKey     string `yaml:"openai_api_key"`
		TelegramBotToken string `yaml:"telegram_bot_token"`
		TelegramChatID   string `yaml:"telegram_chat_id"`
	} `yaml:"credentials"`
	Files struct {
		CredentialsFile string `yaml:"credentials_file"`
		TokenFile       string `yaml:"token_file"`
		PromptsFile     string `yaml:"prompts_file"`
	} `yaml:"files"`
	Polling struct {
		IntervalSeconds int `yaml:"interval_seconds"`
	} `yaml:"polling"`
	OpenAI struct {
		Endpoint    string `yaml:"endpoint"`
		Model       string `yaml:"model"`
		MaxTokens   int    `yaml:"max_tokens"`
		Temperature int    `yaml:"temperature"`
	} `yaml:"openai"`
	Telegram struct {
		ImportantEmailTemplate string `yaml:"important_email_template"`
	} `yaml:"telegram"`
}

// LoadConfig reads and parses the configuration YAML file
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := &Config{}

	if err = yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config YAML: %w", err)
	}

	return config, nil
}

func ValidateConfig(config *Config) error {
	if config.Files.CredentialsFile == "" {
		return fmt.Errorf("credentials_file is required in config.yaml")
	}

	if config.Files.TokenFile == "" {
		return fmt.Errorf("token_file is required in config.yaml")
	}

	if config.Files.PromptsFile == "" {
		return fmt.Errorf("prompts_file is required in config.yaml")
	}

	// Validate credentials from config
	if config.Credentials.OpenAIAPIKey == "" {
		log.Fatal("openai_api_key is required in config.yaml")
	}

	if config.Credentials.TelegramBotToken == "" {
		log.Fatal("telegram_bot_token is required in config.yaml")
	}

	if config.Credentials.TelegramChatID == "" {
		log.Fatal("telegram_chat_id is required in config.yaml")
	}

	if config.Polling.IntervalSeconds <= 0 {
		log.Fatal("interval_seconds must be greater than 0 in config.yaml")
	}

	if config.OpenAI.MaxTokens <= 0 {
		log.Fatal("max_tokens must be greater than 0 in config.yaml")
	}

	if config.OpenAI.Temperature < 0 || config.OpenAI.Temperature > 1 {
		log.Fatal("temperature must be between 0 and 1 in config.yaml")
	}

	if config.OpenAI.Endpoint == "" {
		log.Fatal("endpoint is required in config.yaml")
	}

	if config.Telegram.ImportantEmailTemplate == "" {
		log.Fatal("important_email_template is required in config.yaml")
	}

	return nil
}
