package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var telegramApiBaseUrl = "https://api.telegram.org"

// SendMessage sends a message to a Telegram chat.
func SendMessage(botToken, chatID, text string) error {
	url := fmt.Sprintf("%s/bot%s/sendMessage", telegramApiBaseUrl, botToken)
	payload := map[string]string{
		"chat_id":    chatID,
		"text":       text,
		"parse_mode": "Markdown",
	}

	b, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", strings.NewReader(string(b)))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("telegram API error: %d %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
