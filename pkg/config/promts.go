package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Prompts holds the prompt templates loaded from YAML
type Prompts struct {
	EmailClassification struct {
		SystemMessage      string `yaml:"system_message"`
		UserPromptTemplate string `yaml:"user_prompt_template"`
	} `yaml:"email_classification"`
}

// LoadPrompts reads and parses the prompts YAML file
func LoadPrompts(filename string) (*Prompts, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read prompts file: %w", err)
	}

	prompts := &Prompts{}
	if err := yaml.Unmarshal(data, prompts); err != nil {
		return nil, fmt.Errorf("failed to parse prompts YAML: %w", err)
	}

	return prompts, nil
}
