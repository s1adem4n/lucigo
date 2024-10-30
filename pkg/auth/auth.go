package auth

import (
	"context"
	"time"

	"golang.org/x/oauth2"
)

type OAuthProvider interface {
	GetAuthURL() string
	ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error)
	GetEmail(ctx context.Context, token *oauth2.Token) (string, error)
	GetName() string
}

type Connection interface {
	GetUserID() string
	GetProvider() string
	GetEmail() string
	GetToken() string
	GetRefreshToken() string
}

type User interface {
	GetID() string
}

type Session interface {
	// the id of the session, a hash of the token
	GetID() string
	GetUserID() string
	GetExpiry() *time.Time
}

type Database interface {
	GetConnectionByProviderAndEmail(ctx context.Context, provider, email string) (Connection, error)
	GetUser(ctx context.Context, id string) (User, error)
	GetSession(ctx context.Context, id string) (Session, error)

	CreateUser(ctx context.Context, email string) (User, error)
	CreateConnection(ctx context.Context, userID, provider, email string) (Connection, error)
	CreateSession(ctx context.Context, id, userID string, expiry time.Time) (Session, error)

	DeleteSession(ctx context.Context, id string) error
}

type ProviderMap map[string]OAuthProvider

func (pm ProviderMap) Get(provider string) (OAuthProvider, bool) {
	p, ok := pm[provider]
	return p, ok
}

type Auth struct {
	db        Database
	providers ProviderMap
}

func NewAuth(db Database) *Auth {
	return &Auth{
		db: db,
	}
}

func (a *Auth) RegisterOrLoginOAuth2(ctx context.Context, token string, provider OAuthProvider, code string) (Session, error) {
	t, err := provider.ExchangeCode(ctx, code)
	if err != nil {
		return nil, err
	}

	email, err := provider.GetEmail(ctx, t)
	if err != nil {
		return nil, err
	}

	c, err := a.db.GetConnectionByProviderAndEmail(ctx, provider.GetName(), email)
	if err != nil {
		user, err := a.db.CreateUser(ctx, email)
		if err != nil {
			return nil, err
		}

		c, err = a.db.CreateConnection(ctx, user.GetID(), provider.GetName(), email)
		if err != nil {
			return nil, err
		}
	}

	s, err := a.db.CreateSession(ctx, getSessionID(token), c.GetUserID(), time.Now().Add(24*time.Hour))
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (a *Auth) GetSession(ctx context.Context, token string) (Session, error) {
	return a.db.GetSession(ctx, getSessionID(token))
}

func (a *Auth) DeleteSession(ctx context.Context, token string) error {
	return a.db.DeleteSession(ctx, getSessionID(token))
}
