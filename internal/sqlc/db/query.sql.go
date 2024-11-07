// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package db

import (
	"context"
)

const cleanupSessions = `-- name: CleanupSessions :exec
DELETE FROM Sessions WHERE expires <= ?
`

func (q *Queries) CleanupSessions(ctx context.Context, expires int64) error {
	_, err := q.db.ExecContext(ctx, cleanupSessions, expires)
	return err
}

const createSession = `-- name: CreateSession :one
INSERT INTO Sessions (id, created, updated, expires, data) VALUES (?, ?, ?, ?, ?) RETURNING id, created, updated, expires, data
`

type CreateSessionParams struct {
	ID      []byte
	Created int64
	Updated int64
	Expires int64
	Data    []byte
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error) {
	row := q.db.QueryRowContext(ctx, createSession,
		arg.ID,
		arg.Created,
		arg.Updated,
		arg.Expires,
		arg.Data,
	)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Updated,
		&i.Expires,
		&i.Data,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO Users (id, created, email, name) VALUES (?, ?, ?, ?) RETURNING id, created, email, name
`

type CreateUserParams struct {
	ID      []byte
	Created int64
	Email   string
	Name    string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.Created,
		arg.Email,
		arg.Name,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Email,
		&i.Name,
	)
	return i, err
}

const deleteSession = `-- name: DeleteSession :exec
DELETE FROM Sessions WHERE id = ?
`

func (q *Queries) DeleteSession(ctx context.Context, id []byte) error {
	_, err := q.db.ExecContext(ctx, deleteSession, id)
	return err
}

const getSession = `-- name: GetSession :one
Select id, created, updated, expires, data from Sessions where id = ? LIMIT 1
`

func (q *Queries) GetSession(ctx context.Context, id []byte) (Session, error) {
	row := q.db.QueryRowContext(ctx, getSession, id)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Updated,
		&i.Expires,
		&i.Data,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, created, email, name FROM Users WHERE id = ? LIMIT 1
`

func (q *Queries) GetUserByID(ctx context.Context, id []byte) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Email,
		&i.Name,
	)
	return i, err
}

const updateSessionData = `-- name: UpdateSessionData :exec
UPDATE Sessions SET updated = ?, expires = ?, data = ? WHERE id = ?
`

type UpdateSessionDataParams struct {
	Updated int64
	Expires int64
	Data    []byte
	ID      []byte
}

func (q *Queries) UpdateSessionData(ctx context.Context, arg UpdateSessionDataParams) error {
	_, err := q.db.ExecContext(ctx, updateSessionData,
		arg.Updated,
		arg.Expires,
		arg.Data,
		arg.ID,
	)
	return err
}
