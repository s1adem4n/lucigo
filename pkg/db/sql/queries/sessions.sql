-- name: GetSession :one
SELECT
  *
FROM
  sessions
WHERE
  id = ?;

-- name: CreateSession :one
INSERT INTO
  sessions (id, user_id, expiry)
VALUES
  (?, ?, ?) RETURNING *;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE
  id = ?;