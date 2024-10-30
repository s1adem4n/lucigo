package auth

import (
	"context"
	"encoding/json"
	"errors"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

type GithubOAuth2Provider struct {
	config *oauth2.Config
}

func NewGithubOAuth2Provider(clientID, clientSecret, redirectURL string) OAuthProvider {
	return &GithubOAuth2Provider{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Endpoint:     endpoints.GitHub,
		},
	}
}

func (p *GithubOAuth2Provider) GetName() string {
	return "github"
}

func (p *GithubOAuth2Provider) GetAuthURL() string {
	return p.config.AuthCodeURL("state")
}

func (p *GithubOAuth2Provider) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := p.config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (p *GithubOAuth2Provider) GetEmail(ctx context.Context, token *oauth2.Token) (string, error) {
	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}

	res, err := p.config.Client(ctx, token).Get("https://api.github.com/user/emails")
	if err != nil {
		return "", err
	}

	if err := json.NewDecoder(res.Body).Decode(&emails); err != nil {
		return "", err
	}

	var email string
	for _, e := range emails {
		if e.Primary && e.Verified {
			email = e.Email
			break
		}
	}

	if email == "" {
		return "", errors.New("no verified primary email found")
	}

	return email, nil
}
