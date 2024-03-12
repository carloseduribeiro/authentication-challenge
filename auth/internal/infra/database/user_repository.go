package database

import (
	"context"
	"errors"
	"github.com/carloseduribeiro/auth-challenge/auth/internal/domain/entity/user"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	q *Queries
}

func NewUserRepository(dbPool *pgxpool.Pool) *UserRepository {
	return &UserRepository{q: New(dbPool)}
}

var ErrUserNotFound = errors.New("user not found")

// TODO - testar
func (r *UserRepository) GetUserByDocument(ctx context.Context, document string) (*user.User, error) {
	m, err := r.q.GetUserByDocument(ctx, document)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user.NewUser(m.Document, m.Name, m.Email, m.Password, m.Birthdate, user.WithID(m.ID))
}

// TODO - testar
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	m, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user.NewUser(m.Document, m.Name, m.Email, m.Password, m.Birthdate, user.WithID(m.ID))
}

// TODO - testar
func (r *UserRepository) Create(ctx context.Context, user *user.User) error {
	params := InsertUserParams{
		ID:        user.ID(),
		Document:  user.Document(),
		Name:      user.Name(),
		Email:     user.Email(),
		Password:  user.Password(),
		Birthdate: user.BirthDate(),
		Type:      AuthUserType(user.Type()),
	}
	return r.q.InsertUser(ctx, params)
}
