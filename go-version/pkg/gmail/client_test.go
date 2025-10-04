package gmail

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/oauth2"
)

func TestSaveToken(t *testing.T) {
	t.Run("saves token successfully", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "test")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(dir)

		path := filepath.Join(dir, "token.json")
		token := &oauth2.Token{
			AccessToken:  "access",
			RefreshToken: "refresh",
		}

		err = saveToken(path, token)
		if err != nil {
			t.Fatalf("saveToken() error = %v, wantErr %v", err, false)
		}

		f, err := os.Open(path)
		if err != nil {
			t.Fatalf("Failed to open token file: %v", err)
		}
		defer f.Close()

		var gotToken oauth2.Token
		err = json.NewDecoder(f).Decode(&gotToken)
		if err != nil {
			t.Fatalf("Failed to decode token: %v", err)
		}

		if gotToken.AccessToken != token.AccessToken {
			t.Errorf("AccessToken got = %v, want %v", gotToken.AccessToken, token.AccessToken)
		}
		if gotToken.RefreshToken != token.RefreshToken {
			t.Errorf("RefreshToken got = %v, want %v", gotToken.RefreshToken, token.RefreshToken)
		}
	})

	t.Run("returns error on invalid path", func(t *testing.T) {
		token := &oauth2.Token{
			AccessToken:  "access",
			RefreshToken: "refresh",
		}

		err := saveToken("/non-existent-dir/token.json", token)
		if err == nil {
			t.Error("saveToken() error = nil, wantErr not nil")
		}
	})
}

func TestTokenFromFile(t *testing.T) {
	t.Run("reads token successfully", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "test")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(dir)

		path := filepath.Join(dir, "token.json")
		token := &oauth2.Token{
			AccessToken:  "access",
			RefreshToken: "refresh",
		}

		f, err := os.Create(path)
		if err != nil {
			t.Fatalf("Failed to create token file: %v", err)
		}

		err = json.NewEncoder(f).Encode(token)
		if err != nil {
			f.Close()
			t.Fatalf("Failed to write to token file: %v", err)
		}
		f.Close()

		gotToken, err := tokenFromFile(path)
		if err != nil {
			t.Fatalf("tokenFromFile() error = %v, wantErr %v", err, false)
		}

		if gotToken.AccessToken != token.AccessToken {
			t.Errorf("AccessToken got = %v, want %v", gotToken.AccessToken, token.AccessToken)
		}
		if gotToken.RefreshToken != token.RefreshToken {
			t.Errorf("RefreshToken got = %v, want %v", gotToken.RefreshToken, token.RefreshToken)
		}
	})

	t.Run("returns error on non-existent file", func(t *testing.T) {
		_, err := tokenFromFile("/non-existent-dir/token.json")
		if err == nil {
			t.Error("tokenFromFile() error = nil, wantErr not nil")
		}
	})
}
