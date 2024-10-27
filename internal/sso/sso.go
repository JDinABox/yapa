package sso

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type Provider struct {
	Provider *oidc.Provider
	Config   *oauth2.Config
	Verifier *oidc.IDTokenVerifier
}

type Providers struct {
	Apple  *Provider
	Google *Provider
}

type Claims struct {
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
}

func GetScopes(additionalScopes ...string) []string {
	var scopes []string
	scopes = append(scopes, oidc.ScopeOpenID)
	scopes = append(scopes, "email")
	scopes = append(scopes, additionalScopes...)
	return scopes
}

func NewProvider(clientID, clientSecret, providerURL, redirectURL string, scopes []string) (*Provider, error) {
	var err error
	ctx := context.Background()
	p := &Provider{}
	p.Provider, err = oidc.NewProvider(ctx, providerURL)
	if err != nil {
		return nil, err
	}
	p.Config = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		Endpoint:     p.Provider.Endpoint(),
	}

	p.Verifier = p.Provider.Verifier(&oidc.Config{ClientID: clientID})

	return p, nil
}
