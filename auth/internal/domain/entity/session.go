package entity

import (
	"github.com/google/uuid"
	"time"
)

type Session struct {
	id        uuid.UUID
	userId    uuid.UUID
	createdAt time.Time
	expiresAt time.Time
}

func NewSession(id, userId uuid.UUID, createdAt time.Time, d time.Duration) *Session {
	return &Session{
		id:        id,
		userId:    userId,
		createdAt: createdAt,
		expiresAt: createdAt.Add(d),
	}
}

func (s Session) Id() uuid.UUID {
	return s.id
}
func (s Session) UserId() uuid.UUID {
	return s.userId
}
func (s Session) CreatedAt() time.Time {
	return s.createdAt
}
func (s Session) ExpiresAt() time.Time {
	return s.expiresAt
}
