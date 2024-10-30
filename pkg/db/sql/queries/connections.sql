-- name: GetConnectionByProviderAndEmail :one
SELECT
  *
FROM
  connections
WHERE
  provider = ? AND
  email = ?;

-- name: CreateConnection :one
INSERT INTO
  connections (
    id,
    user_id,
    provider,
    email,
    token,
    refresh_token
  )
VALUES
  (?, ?, ?, ?, ?, ?) RETURNING *;