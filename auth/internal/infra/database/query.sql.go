// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const getUserByDocument = `-- name: GetUserByDocument :one
SELECT id, document, name, email, password, birthdate, type
FROM auth.users
WHERE document = $1
`

func (q *Queries) GetUserByDocument(ctx context.Context, document string) (AuthUser, error) {
	row := q.db.QueryRow(ctx, getUserByDocument, document)
	var i AuthUser
	err := row.Scan(
		&i.ID,
		&i.Document,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.Birthdate,
		&i.Type,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, document, name, email, password, birthdate, type
FROM auth.users
WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (AuthUser, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i AuthUser
	err := row.Scan(
		&i.ID,
		&i.Document,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.Birthdate,
		&i.Type,
	)
	return i, err
}

const insertUser = `-- name: InsertUser :exec
INSERT INTO auth.users (id, document, name, email, password, birthdate, type)
VALUES ($1, $2, $3, $4, $5, $6, $7)
`

type InsertUserParams struct {
	ID        uuid.UUID
	Document  string
	Name      string
	Email     string
	Password  string
	Birthdate time.Time
	Type      AuthUserType
}

func (q *Queries) InsertUser(ctx context.Context, arg InsertUserParams) error {
	_, err := q.db.Exec(ctx, insertUser,
		arg.ID,
		arg.Document,
		arg.Name,
		arg.Email,
		arg.Password,
		arg.Birthdate,
		arg.Type,
	)
	return err
}
