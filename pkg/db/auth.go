package db

import (
	"context"
	"lucigo/pkg/auth"
	"time"

	"github.com/oklog/ulid/v2"
)

func generateID() string {
	id, _ := ulid.New(ulid.Timestamp(time.Now()), nil)
	return id.String()
}

func (c Connection) GetUserID() string {
	return c.UserID
}
func (c Connection) GetProvider() string {
	return c.Provider
}
func (c Connection) GetEmail() string {
	return c.Email
}
func (c Connection) GetToken() string {
	return c.Token
}
func (c Connection) GetRefreshToken() string {
	return c.RefreshToken.String
}

func (u User) GetID() string {
	return u.ID
}

func (s Session) GetID() string {
	return s.ID
}
func (s Session) GetUserID() string {
	return s.UserID
}
func (s Session) GetExpiry() *time.Time {
	t := time.Unix(s.Expiry, 0)
	return &t
}

type AuthDatabase struct {
	queries *Queries
}

func NewAuthDatabase(queries *Queries) auth.Database {
	return &AuthDatabase{
		queries: queries,
	}
}

func (a *AuthDatabase) GetConnectionByProviderAndEmail(ctx context.Context, provider, email string) (auth.Connection, error) {
	return a.queries.GetConnectionByProviderAndEmail(ctx, GetConnectionByProviderAndEmailParams{
		Provider: provider,
		Email:    email,
	})
}
func (a *AuthDatabase) GetUser(ctx context.Context, id string) (auth.User, error) {
	return a.queries.GetUser(ctx, id)
}
func (a *AuthDatabase) GetSession(ctx context.Context, id string) (auth.Session, error) {
	return a.queries.GetSession(ctx, id)
}

func (a *AuthDatabase) CreateUser(ctx context.Context, email string) (auth.User, error) {
	return a.queries.CreateUser(ctx, CreateUserParams{
		ID:     generateID(),
		Active: true,
	})
}
func (a *AuthDatabase) CreateConnection(ctx context.Context, userID, provider, email string) (auth.Connection, error) {
	return a.queries.CreateConnection(ctx, CreateConnectionParams{
		ID:       generateID(),
		UserID:   userID,
		Provider: provider,
		Email:    email,
		Token:    "",
	})
}
func (a *AuthDatabase) CreateSession(ctx context.Context, id, userID string, expiry time.Time) (auth.Session, error) {
	return a.queries.CreateSession(ctx, CreateSessionParams{
		ID:     id,
		UserID: userID,
		Expiry: expiry.Unix(),
	})
}

func (a *AuthDatabase) DeleteSession(ctx context.Context, id string) error {
	return a.queries.DeleteSession(ctx, id)
}
