package sso_test

import (
	"testing"

	openid "github.com/JDinABox/yapa/internal/sso"
	"golang.org/x/oauth2"
)

func TestGetScopes(t *testing.T) {
	scopes := openid.GetScopes("profile")
	if len(scopes) != 3 {
		t.Errorf("expected 3 scopes, got %d", len(scopes))
	}
}

func TestNewProvider(t *testing.T) {
	clientID := "com.example.client"
	clientSecret := "example-client-secret"
	providerURL := "https://appleid.apple.com"
	redirectURL := "https://example.com/callback"
	scopes := []string{"openid", "email"}

	provider, err := openid.NewProvider(clientID, clientSecret, providerURL, redirectURL, scopes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if provider.Provider == nil {
		t.Errorf("expected Provider field to be set, got nil")
	}

	if provider.Config == nil {
		t.Errorf("expected Config field to be set, got nil")
	}

	if provider.Verifier == nil {
		t.Errorf("expected Verifier field to be set, got nil")
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		Endpoint:     provider.Provider.Endpoint(),
	}
	if provider.Config.ClientID != config.ClientID ||
		provider.Config.ClientSecret != config.ClientSecret ||
		provider.Config.RedirectURL != config.RedirectURL ||
		provider.Config.Scopes[0] != config.Scopes[0] ||
		provider.Config.Scopes[1] != config.Scopes[1] ||
		provider.Config.Endpoint != config.Endpoint {
		t.Errorf("provider.Config does not match expected Config")
	}
}
