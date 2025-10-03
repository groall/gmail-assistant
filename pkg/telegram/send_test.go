package telegram

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendMessage(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/botmy-token/sendMessage" {
				t.Errorf("expected path /botmy-token/sendMessage, got %s", r.URL.Path)
			}
			var payload map[string]string
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatal(err)
			}
			if payload["chat_id"] != "my-chat" {
				t.Errorf("expected chat_id my-chat, got %s", payload["chat_id"])
			}
			if payload["text"] != "my-text" {
				t.Errorf("expected text my-text, got %s", payload["text"])
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		telegramApiBaseUrl = server.URL

		err := SendMessage("my-token", "my-chat", "my-text")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("api error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("error from api"))
		}))
		defer server.Close()

		telegramApiBaseUrl = server.URL

		err := SendMessage("my-token", "my-chat", "my-text")
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		expectedError := "telegram API error: 400 error from api"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("http error", func(t *testing.T) {
		telegramApiBaseUrl = ""

		err := SendMessage("my-token", "my-chat", "my-text")
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
	})
}
