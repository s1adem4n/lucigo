// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO
  users (id, active)
VALUES
  (?, ?) RETURNING id, active
`

type CreateUserParams struct {
	ID     string `json:"id"`
	Active bool   `json:"active"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (*User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.ID, arg.Active)
	var i User
	err := row.Scan(&i.ID, &i.Active)
	return &i, err
}

const getUser = `-- name: GetUser :one
SELECT
  id, active
FROM
  users
WHERE
  id = ?
`

func (q *Queries) GetUser(ctx context.Context, id string) (*User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(&i.ID, &i.Active)
	return &i, err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users
SET
  active = ?
WHERE
  id = ?
`

type UpdateUserParams struct {
	Active bool   `json:"active"`
	ID     string `json:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.ExecContext(ctx, updateUser, arg.Active, arg.ID)
	return err
}