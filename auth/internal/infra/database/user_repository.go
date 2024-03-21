package database

import (
	"context"
	"errors"
	"github.com/carloseduribeiro/authentication-challenge/auth/internal/domain/entity"
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

func (r *UserRepository) GetUserByDocument(ctx context.Context, document string) (*entity.User, error) {
	m, err := r.q.GetUserByDocument(ctx, document)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return entity.NewUser(m.ID, m.Document, m.Name, m.Email, m.Birthdate, entity.WithPasswordHashed(m.Password))
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	m, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return entity.NewUser(m.ID, m.Document, m.Name, m.Email, m.Birthdate, entity.WithPasswordHashed(m.Password))
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
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
