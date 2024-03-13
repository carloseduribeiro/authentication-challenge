-- name: GetUserByDocument :one
SELECT *
FROM auth.users
WHERE document = $1;

-- name: GetUserByEmail :one
SELECT *
FROM auth.users
WHERE email = $1;

-- name: InsertUser :exec
INSERT INTO auth.users (id, document, name, email, password, birthdate, type)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: InsertSession :exec
INSERT INTO auth.sessions (id, user_id, created_at, expires_at)
VALUES ($1, $2, $3, $4);