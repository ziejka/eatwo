// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package db

import (
	"context"
)

const createDream = `-- name: CreateDream :one
INSERT INTO
  dreams (id, user_id, description, explanation, date)
VALUES
  (?, ?, ?, ?, ?) RETURNING id, user_id, description, explanation, date
`

type CreateDreamParams struct {
	ID          string
	UserID      string
	Description string
	Explanation string
	Date        string
}

func (q *Queries) CreateDream(ctx context.Context, arg CreateDreamParams) (Dream, error) {
	row := q.db.QueryRowContext(ctx, createDream,
		arg.ID,
		arg.UserID,
		arg.Description,
		arg.Explanation,
		arg.Date,
	)
	var i Dream
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Description,
		&i.Explanation,
		&i.Date,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO
  users (id, email, name, hash_password)
VALUES
  (?, ?, ?, ?) RETURNING id, email, name, hash_password
`

type CreateUserParams struct {
	ID           string
	Email        string
	Name         string
	HashPassword string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.Email,
		arg.Name,
		arg.HashPassword,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Name,
		&i.HashPassword,
	)
	return i, err
}

const getDreams = `-- name: GetDreams :many
SELECT
  id, user_id, description, explanation, date
FROM
  dreams
WHERE
  user_id = ?
ORDER BY
  date
`

func (q *Queries) GetDreams(ctx context.Context, userID string) ([]Dream, error) {
	rows, err := q.db.QueryContext(ctx, getDreams, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Dream
	for rows.Next() {
		var i Dream
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Description,
			&i.Explanation,
			&i.Date,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUser = `-- name: GetUser :one
SELECT
  id, email, name, hash_password
FROM
  users
WHERE
  email = ?
LIMIT
  1
`

func (q *Queries) GetUser(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Name,
		&i.HashPassword,
	)
	return i, err
}

const updateDreamExplanation = `-- name: UpdateDreamExplanation :one
UPDATE dreams
SET
  explanation = ?
WHERE
  id = ?
  AND user_id = ? RETURNING id, user_id, description, explanation, date
`

type UpdateDreamExplanationParams struct {
	Explanation string
	ID          string
	UserID      string
}

func (q *Queries) UpdateDreamExplanation(ctx context.Context, arg UpdateDreamExplanationParams) (Dream, error) {
	row := q.db.QueryRowContext(ctx, updateDreamExplanation, arg.Explanation, arg.ID, arg.UserID)
	var i Dream
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Description,
		&i.Explanation,
		&i.Date,
	)
	return i, err
}
