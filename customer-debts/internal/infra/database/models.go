// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package database

import (
	"time"

	"github.com/google/uuid"
)

type CustomerDebt struct {
	ID        uuid.UUID
	Document  string
	Duedate   time.Time
	Amount    float64
	CreatedAt time.Time
}