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
