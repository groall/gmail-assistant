package classifier

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

type Config struct {
	OpenAI              OpenAIConfig
	EmailClassification EmailClassificationConfig
}

type OpenAIConfig struct {
	Endpoint    string
	Model       string
	MaxTokens   int
	Temperature int
	APIKey      string
}

type EmailClassificationConfig struct {
	SystemMessage      string
	UserPromptTemplate string
}

type Classifier struct {
	openAIClient *openai.Client
	config       *Config
	ctx          context.Context
}

type llmDecision struct {
	Important   bool   `json:"important"`
	Explanation string `json:"explanation"`
}

func NewClassifier(cfg *Config, ctx context.Context) *Classifier {
	c := openai.DefaultConfig(cfg.OpenAI.APIKey)
	c.BaseURL = normalizeOpenAIBaseURL(cfg.OpenAI.Endpoint)
	client := openai.NewClientWithConfig(c)
	return &Classifier{
		openAIClient: client,
		config:       cfg,
		ctx:          ctx,
	}

}

// normalizeOpenAIBaseURL ensures the BaseURL is suitable for go-openai client
// - appends "/v1" if missing
// - trims any path after "/v1/" if a full endpoint URL was provided
func normalizeOpenAIBaseURL(endpoint string) string {
	if endpoint == "" {
		return ""
	}

	e := strings.TrimRight(endpoint, "/")
	if strings.HasSuffix(e, "/v1") {
		return e
	}

	if idx := strings.Index(e, "/v1/"); idx != -1 {
		return e[:idx+3]
	}

	return e + "/v1"
}

// ClassifyEmail calls OpenAI Chat Completion API to decide if an email is important.
// Returns (important bool, explanation string).
func (c *Classifier) ClassifyEmail(text string) (bool, string) {
	// Build a concise prompt and send to OpenAI-compatible API using SDK
	prompt := fmt.Sprintf(c.config.EmailClassification.UserPromptTemplate, text)

	ctx, cancel := context.WithTimeout(c.ctx, 30*time.Second)
	defer cancel()

	resp, err := c.openAIClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       c.config.OpenAI.Model,
		MaxTokens:   c.config.OpenAI.MaxTokens,
		Temperature: float32(c.config.OpenAI.Temperature),
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: c.config.EmailClassification.SystemMessage},
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	})
	if err != nil {
		log.Fatalf("OpenAI request failed: %v", err)
		// return false, "request failed"
	}
	if len(resp.Choices) == 0 {
		log.Printf("No choices in OpenAI response")
		return false, "no choices"
	}

	content := strings.TrimSpace(resp.Choices[0].Message.Content)

	return parseLLMDecision(content)
}

func parseLLMDecision(content string) (bool, string) {
	// Try to extract JSON from the assistant content
	// find first '{' and last '}'
	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")
	if start != -1 && end != -1 && end > start {
		jsonStr := content[start : end+1]
		var j llmDecision
		if err := json.Unmarshal([]byte(jsonStr), &j); err == nil {
			return j.Important, j.Explanation
		}
	}

	// Fallback: simple heuristic, check for keywords
	lower := strings.ToLower(content)
	if strings.Contains(lower, "true") || strings.Contains(lower, "important") || strings.Contains(lower, "yes") {
		return true, strings.SplitN(content, "\n", 2)[0]
	}

	return false, strings.SplitN(content, "\n", 2)[0]
}
