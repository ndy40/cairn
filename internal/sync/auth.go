package sync

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"golang.org/x/oauth2"
)

var dropboxOAuthEndpoint = oauth2.Endpoint{
	AuthURL:  "https://www.dropbox.com/oauth2/authorize",
	TokenURL: "https://api.dropboxapi.com/oauth2/token",
}

// RunOAuth2Flow runs the Dropbox OAuth2 PKCE flow with no-redirect (manual code entry).
// It prints the authorization URL to stdout and reads the authorization code from stdin.
func RunOAuth2Flow(appKey string) (*oauth2.Token, error) {
	cfg := &oauth2.Config{
		ClientID: appKey,
		Endpoint: dropboxOAuthEndpoint,
		Scopes: []string{
			"files.metadata.read",
			"files.content.read",
			"files.content.write",
		},
	}

	// Generate PKCE code verifier and challenge.
	verifier, err := generateCodeVerifier()
	if err != nil {
		return nil, fmt.Errorf("generate code verifier: %w", err)
	}
	challenge := generateCodeChallenge(verifier)

	// Build authorization URL.
	authURL := cfg.AuthCodeURL("",
		oauth2.SetAuthURLParam("response_type", "code"),
		oauth2.SetAuthURLParam("code_challenge", challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("token_access_type", "offline"),
	)

	fmt.Println("Open this URL in your browser to authorize cairn:")
	fmt.Println()
	fmt.Println(authURL)
	fmt.Println()
	fmt.Print("Enter the authorization code: ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return nil, fmt.Errorf("failed to read authorization code")
	}
	code := strings.TrimSpace(scanner.Text())
	if code == "" {
		return nil, fmt.Errorf("authorization code is empty")
	}

	// Exchange code for token.
	token, err := cfg.Exchange(context.Background(), code,
		oauth2.SetAuthURLParam("code_verifier", verifier),
	)
	if err != nil {
		return nil, fmt.Errorf("exchange code for token: %w", err)
	}

	return token, nil
}

// generateCodeVerifier creates a random 32-byte code verifier for PKCE.
func generateCodeVerifier() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// generateCodeChallenge computes the S256 code challenge from a verifier.
func generateCodeChallenge(verifier string) string {
	h := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h[:])
}
