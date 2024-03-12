package user

import (
	"context"
)

type Repository interface {
	GetUserByDocument(ctx context.Context, document string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error
}
