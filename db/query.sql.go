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

const createItem = `-- name: CreateItem :one
INSERT INTO
    items (value, list_id)
VALUES
    (?, ?) RETURNING id, value, list_id
`

type CreateItemParams struct {
	Value  string
	ListID int64
}

func (q *Queries) CreateItem(ctx context.Context, arg CreateItemParams) (Item, error) {
	row := q.db.QueryRowContext(ctx, createItem, arg.Value, arg.ListID)
	var i Item
	err := row.Scan(&i.ID, &i.Value, &i.ListID)
	return i, err
}

const createList = `-- name: CreateList :one
INSERT INTO
    lists (name, user_id)
VALUES
    (?, ?) RETURNING id, name, user_id
`

type CreateListParams struct {
	Name   string
	UserID string
}

func (q *Queries) CreateList(ctx context.Context, arg CreateListParams) (List, error) {
	row := q.db.QueryRowContext(ctx, createList, arg.Name, arg.UserID)
	var i List
	err := row.Scan(&i.ID, &i.Name, &i.UserID)
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

const deleteDreamsForUser = `-- name: DeleteDreamsForUser :exec
DELETE FROM dreams
WHERE
    user_id = ?
`

func (q *Queries) DeleteDreamsForUser(ctx context.Context, userID string) error {
	_, err := q.db.ExecContext(ctx, deleteDreamsForUser, userID)
	return err
}

const deleteItemsForUser = `-- name: DeleteItemsForUser :exec
DELETE FROM items
WHERE
    list_id IN (
        SELECT
            id
        FROM
            lists
        WHERE
            user_id = ?
    )
`

func (q *Queries) DeleteItemsForUser(ctx context.Context, userID string) error {
	_, err := q.db.ExecContext(ctx, deleteItemsForUser, userID)
	return err
}

const deleteListsForUser = `-- name: DeleteListsForUser :exec
DELETE FROM lists
WHERE
    user_id = ?
`

func (q *Queries) DeleteListsForUser(ctx context.Context, userID string) error {
	_, err := q.db.ExecContext(ctx, deleteListsForUser, userID)
	return err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE
    id = ?
`

func (q *Queries) DeleteUser(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getCheckListByUser = `-- name: GetCheckListByUser :many
SELECT
    l.id,
    l.user_id,
    l.name
FROM
    lists l
WHERE
    l.user_id = ?
`

type GetCheckListByUserRow struct {
	ID     int64
	UserID string
	Name   string
}

func (q *Queries) GetCheckListByUser(ctx context.Context, userID string) ([]GetCheckListByUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getCheckListByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCheckListByUserRow
	for rows.Next() {
		var i GetCheckListByUserRow
		if err := rows.Scan(&i.ID, &i.UserID, &i.Name); err != nil {
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

const getDreams = `-- name: GetDreams :many
SELECT
    id, user_id, description, explanation, date
FROM
    dreams
WHERE
    user_id = ?
ORDER BY
    date DESC
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

const getListIDByUser = `-- name: GetListIDByUser :many
SELECT
    id
FROM
    lists
WHERE
    id = ?
    AND user_id = ?
`

type GetListIDByUserParams struct {
	ID     int64
	UserID string
}

func (q *Queries) GetListIDByUser(ctx context.Context, arg GetListIDByUserParams) ([]int64, error) {
	rows, err := q.db.QueryContext(ctx, getListIDByUser, arg.ID, arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getListWithItemsByListId = `-- name: GetListWithItemsByListId :many
SELECT
    l.id,
    l.name,
    i.id,
    i."value"
FROM
    lists l
    JOIN items i ON i.list_id = l.id
WHERE
    l.id = ?
    AND l.user_id = ?
`

type GetListWithItemsByListIdParams struct {
	ID     int64
	UserID string
}

type GetListWithItemsByListIdRow struct {
	ID    int64
	Name  string
	ID_2  int64
	Value string
}

func (q *Queries) GetListWithItemsByListId(ctx context.Context, arg GetListWithItemsByListIdParams) ([]GetListWithItemsByListIdRow, error) {
	rows, err := q.db.QueryContext(ctx, getListWithItemsByListId, arg.ID, arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetListWithItemsByListIdRow
	for rows.Next() {
		var i GetListWithItemsByListIdRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.ID_2,
			&i.Value,
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

const getUserByID = `-- name: GetUserByID :one
SELECT
    id, email, name, hash_password
FROM
    users
WHERE
    id = ?
LIMIT
    1
`

func (q *Queries) GetUserByID(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, id)
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

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET
    email = ?,
    name = ?
WHERE
    id = ? RETURNING id, email, name, hash_password
`

type UpdateUserParams struct {
	Email string
	Name  string
	ID    string
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser, arg.Email, arg.Name, arg.ID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Name,
		&i.HashPassword,
	)
	return i, err
}
