// Gmail → AI → Telegram agent (Go)
// --------------------------------------------------------
// Single-file prototype that:
//  - uses Gmail API (OAuth2) to fetch unread emails
//  - classifies them via OpenAI (chat completion)
//  - trashes unimportant emails, notifies you on Telegram about important ones
//  - marks important emails as read
//
// Requirements & setup (summary):
// 1) Enable Gmail API in Google Cloud Console and create OAuth 2.0 Client ID (Desktop or Web).
//    Download credentials.json and place in the same folder as this program.
// 2) Configure config.yaml with your API credentials:
//       openai_api_key — your OpenAI API key
//       telegram_bot_token — your Telegram bot token (BotFather)
//       telegram_chat_id — chat ID to receive messages (your user id or group id)
// 3) Run once to get token.json (the OAuth flow will open a browser). The program will save token.json.
// 4) go run main.go
//
// Note: This is a prototype. For production you should:
//  - persist per-sender rules and allow user feedback
//  - handle quota/backoff and exponential retries
//  - secure credentials and token storage
//  - run as a service (Docker, systemd) and/or use Gmail push notifications
//
// --------------------------------------------------------

package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"google.golang.org/api/gmail/v1"

	"gmail-local-agent/pkg/classifier"
	cmdConfig "gmail-local-agent/pkg/config"
	gmail2 "gmail-local-agent/pkg/gmail"
	"gmail-local-agent/pkg/telegram"
)

const configFile = "configs/config.yaml" // main configuration file

var config *cmdConfig.Config

func main() {
	ctx := context.Background()

	var err error
	// Load configuration from YAML
	if config, err = cmdConfig.LoadConfig(configFile); err != nil {
		log.Fatalf("Unable to load config from %s: %v", configFile, err)
	}

	if err = cmdConfig.ValidateConfig(config); err != nil {
		log.Fatalf("Invalid config: %v", err)
	}

	var prompts *cmdConfig.Prompts
	// Load prompts from DB
	if prompts, err = cmdConfig.LoadPrompts(config.Files.PromptsFile); err != nil {
		log.Fatalf("Unable to load prompts from DB: %v", err)
	}

	gmailSrv, err := gmail2.NewService(ctx, config.Files.CredentialsFile, config.Files.TokenFile)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	log.Println("Agent started — polling Gmail for unread messages...")

	pollInterval := time.Duration(config.Polling.IntervalSeconds) * time.Second
	// Create classifier
	clr := newClassifier(ctx, config, prompts)
	for {
		processInbox(clr, gmailSrv)
		time.Sleep(pollInterval)
	}
}

// newClassifier creates a new classifier instance.
func newClassifier(ctx context.Context, config *cmdConfig.Config, prompts *cmdConfig.Prompts) *classifier.Classifier {
	classifierConfig := &classifier.Config{
		OpenAI: classifier.OpenAIConfig{
			APIKey:      config.Credentials.OpenAIAPIKey,
			Model:       config.OpenAI.Model,
			MaxTokens:   config.OpenAI.MaxTokens,
			Temperature: config.OpenAI.Temperature,
			Endpoint:    config.OpenAI.Endpoint,
		},
		EmailClassification: classifier.EmailClassificationConfig{
			SystemMessage:      prompts.EmailClassification.SystemMessage,
			UserPromptTemplate: prompts.EmailClassification.UserPromptTemplate,
		},
	}

	return classifier.NewClassifier(classifierConfig, ctx)
}

// processInbox processes all unread messages in INBOX
func processInbox(clr *classifier.Classifier, srv *gmail.Service) {
	// List unread messages in INBOX
	req := srv.Users.Messages.List("me").Q("is:unread in:inbox")
	resp, err := req.Do()
	if err != nil {
		log.Printf("Unable to fetch messages: %v", err)
		return
	}
	if resp.ResultSizeEstimate == 0 || len(resp.Messages) == 0 {
		log.Println("No unread messages.")
		return
	}

	log.Printf("Processing %d unread messages...", len(resp.Messages))

	for _, m := range resp.Messages {
		processEmail(clr, srv, config.Credentials.TelegramBotToken, config.Credentials.TelegramChatID, m)
	}
}

// processEmail processes a single email
func processEmail(clr *classifier.Classifier, srv *gmail.Service, telegramToken, chatID string, msg *gmail.Message) {
	msgFull, err := srv.Users.Messages.Get("me", msg.Id).Format("full").Do()
	if err != nil {
		log.Printf("Could not retrieve message %s: %v", msg.Id, err)
		return
	}

	snippet := msgFull.Snippet
	subject := getHeader(msgFull.Payload.Headers, "Subject")
	from := getHeader(msgFull.Payload.Headers, "From")
	preview := snippet

	// Compose text for classifier
	inputText := fmt.Sprintf("From: %s\nSubject: %s\n\n%s", from, subject, preview)

	important, reason := clr.ClassifyEmail(inputText)
	if important {
		// Send Telegram notification, mark as read
		body := fmt.Sprintf(config.Telegram.ImportantEmailTemplate, escapeMarkdown(from), escapeMarkdown(subject), escapeMarkdown(preview), escapeMarkdown(reason))
		err = telegram.SendMessage(telegramToken, chatID, body)
		if err != nil {
			log.Printf("Failed to send Telegram message: %v", err)
		} else {
			// Mark as read to avoid re-processing
			_, err = srv.Users.Messages.Modify("me", msg.Id, &gmail.ModifyMessageRequest{RemoveLabelIds: []string{"UNREAD"}}).Do()
			if err != nil {
				log.Printf("Failed to mark message read: %v", err)
			}
		}
	} else {
		msgAboutTrashed := fmt.Sprintf("🗑 Trashed message from %s subject=%s", from, subject)
		err = telegram.SendMessage(telegramToken, chatID, msgAboutTrashed)
		if err != nil {
			log.Printf("Failed to send Telegram message: %v", err)
		} else {
			// Trash the message as unimportant
			_, err = srv.Users.Messages.Trash("me", msg.Id).Do()
			if err != nil {
				log.Printf("Failed to trash message %s: %v", msg.Id, err)
			} else {
				log.Println(msgAboutTrashed)
			}
		}
	}
}

func getHeader(headers []*gmail.MessagePartHeader, name string) string {
	for _, h := range headers {
		if strings.EqualFold(h.Name, name) {
			return h.Value
		}
	}
	return ""
}

// simple markdown escape for a few characters
func escapeMarkdown(s string) string {
	replacer := strings.NewReplacer("_", "\\_", "*", "\\*", "[", "\\[", "]", "\\]")
	return replacer.Replace(s)
}
