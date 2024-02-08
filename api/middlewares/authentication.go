package middlewares

import (
	"context"
	"errors"
)

type User struct {
	ID    int
	Name  string
	Email string
}

var ErrInvalidToken = errors.New("invalid token")

func VerifyToken(token string) (*User, error) {

	if token != "valid_token" {
		return nil, ErrInvalidToken
	}

	user := &User{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	return user, nil
}

func WithUserContext(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func GetUserFromContext(ctx context.Context) *User {
	if user, ok := ctx.Value(userContextKey).(*User); ok {
		return user
	}
	return nil
}

var userContextKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func (c *contextKey) String() string {
	return c.name
}
