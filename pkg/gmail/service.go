package gmail

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// NewService creates a new Gmail service.
func NewService(ctx context.Context, credentialsFile string, tokenFile string) (*gmail.Service, error) {
	b, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, fmt.Errorf("Unable to read %s: %w\nFollow setup steps in the repository header.", credentialsFile, err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	gmailConfig, err := google.ConfigFromJSON(b, gmail.GmailModifyScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %w", err)
	}

	// Retrieve a token, saves the token, then returns the generated client.
	gmailClient, err := getClient(ctx, gmailConfig, tokenFile)
	if err != nil {
		return nil, fmt.Errorf("unable to get Gmail client: %w", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(gmailClient))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Gmail client: %w", err)
	}

	return srv, nil
}
