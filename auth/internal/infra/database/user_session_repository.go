package database

import (
	"context"
	"github.com/carloseduribeiro/auth-challenge/auth/internal/domain/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository struct {
	q *Queries
}

func NewSessionRepository(dbPool *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{q: New(dbPool)}
}

func (s *SessionRepository) Save(ctx context.Context, session *entity.Session) error {
	params := InsertSessionParams{
		ID:        session.Id(),
		UserID:    session.UserId(),
		CreatedAt: session.CreatedAt(),
		ExpiresAt: session.ExpiresAt(),
	}
	return s.q.InsertSession(ctx, params)
}
