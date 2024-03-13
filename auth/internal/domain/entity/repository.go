package entity

import (
	"context"
)

type UserRepository interface {
	GetUserByDocument(ctx context.Context, document string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error
}

type SessionRepository interface {
	Save(ctx context.Context, session *Session) error
}
