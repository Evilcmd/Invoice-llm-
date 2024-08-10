// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: users.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :exec
INSERT INTO users (id, user_name, passwd)
VALUES ($1, $2, $3)
`

type CreateUserParams struct {
	ID       uuid.UUID
	UserName string
	Passwd   string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.ExecContext(ctx, createUser, arg.ID, arg.UserName, arg.Passwd)
	return err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, user_name, passwd FROM users WHERE user_name = $1
`

func (q *Queries) GetUserByUsername(ctx context.Context, userName string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByUsername, userName)
	var i User
	err := row.Scan(&i.ID, &i.UserName, &i.Passwd)
	return i, err
}
