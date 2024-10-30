-- name: GetUser :one
SELECT
  *
FROM
  users
WHERE
  id = ?;

-- name: CreateUser :one
INSERT INTO
  users (id, active)
VALUES
  (?, ?) RETURNING *;

-- name: UpdateUser :exec
UPDATE users
SET
  active = ?
WHERE
  id = ?;