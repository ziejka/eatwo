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

-- name: CreateDream :one
INSERT INTO
  dreams (id, user_id, description, explanation, date)
VALUES
  (?, ?, ?, ?, ?) RETURNING *;

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
  date;
