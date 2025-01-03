-- name: GetUser :one
SELECT
    *
FROM
    users
WHERE
    email = ?
LIMIT
    1;

-- name: CreateUser :one
INSERT INTO
    users (id, email, name, hash_password)
VALUES
    (?, ?, ?, ?) RETURNING *;

-- name: GetUserByID :one
SELECT
    *
FROM
    users
WHERE
    id = ?
LIMIT
    1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE
    id = ?;

-- name: UpdateUser :one
UPDATE users
SET
    email = ?,
    name = ?
WHERE
    id = ? RETURNING *;

-- name: CreateDream :one
INSERT INTO
    dreams (id, user_id, description, explanation, date)
VALUES
    (?, ?, ?, ?, ?) RETURNING *;

-- name: DeleteDreamsForUser :exec
DELETE FROM dreams
WHERE
    user_id = ?;

-- name: UpdateDreamExplanation :one
UPDATE dreams
SET
    explanation = ?
WHERE
    id = ?
    AND user_id = ? RETURNING *;

-- name: GetDreams :many
SELECT
    *
FROM
    dreams
WHERE
    user_id = ?
ORDER BY
    date DESC;

-- name: GetCheckListByUser :many
SELECT
    l.id,
    l.user_id,
    l.name
FROM
    lists l
WHERE
    l.user_id = ?;

-- name: GetListWithItemsByListId :many
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
    AND l.user_id = ?;

-- name: DeleteListsForUser :exec
DELETE FROM lists
WHERE
    user_id = ?;

-- name: CreateList :one
INSERT INTO
    lists (name, user_id)
VALUES
    (?, ?) RETURNING *;

-- name: GetListIDByUser :many
SELECT
    id
FROM
    lists
WHERE
    id = ?
    AND user_id = ?;

-- name: CreateItem :one
INSERT INTO
    items (value, list_id)
VALUES
    (?, ?) RETURNING *;

-- name: DeleteItemsForUser :exec
DELETE FROM items
WHERE
    list_id IN (
        SELECT
            id
        FROM
            lists
        WHERE
            user_id = ?
    );
