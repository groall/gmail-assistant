package gmail

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It caches the token in tokenFile.
func getClient(ctx context.Context, gmailConfig *oauth2.Config, tokenFile string) (*http.Client, error) {
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		tok, err = getTokenFromWeb(gmailConfig)
		if err != nil {
			return nil, fmt.Errorf("unable to get token from web: %w", err)
		}

		if err = saveToken(tokenFile, tok); err != nil {
			return nil, fmt.Errorf("unable to cache oauth token: %w", err)
		}
	}

	return gmailConfig.Client(ctx, tok), nil
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(gmailConfig *oauth2.Config) (*oauth2.Token, error) {
	authURL := gmailConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code:\n%v\n", authURL)

	// For convenience try to open browser, but still ask user to paste code
	// (os/exec omitted for cross-platform simplicity)

	fmt.Print("Enter authorization code: ")
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		return nil, fmt.Errorf("unable to read authorization code: %w", err)
	}

	tok, err := gmailConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web: %w", err)
	}

	return tok, nil
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var tok oauth2.Token
	err = json.NewDecoder(f).Decode(&tok)
	return &tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to cache oauth token: %w", err)
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		return fmt.Errorf("unable to cache oauth token: %w", err)
	}

	return nil
}
