-- name: CreateTrader :one
INSERT INTO traders (
  first_name, last_name, username, password, email, country, phone
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetTrader :one
SELECT * FROM traders 
WHERE username = $1
LIMIT 1;

-- name: GetTraders :many
SELECT * FROM traders 
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateTrader :one
UPDATE traders
SET
  first_name = COALESCE(sqlc.narg(first_name), first_name),
  last_name = COALESCE(sqlc.narg(last_name), last_name),
  username = COALESCE(sqlc.narg(username), username),
  password = COALESCE(sqlc.narg(password), password),
  email = COALESCE(sqlc.narg(email), email),
  phone = COALESCE(sqlc.narg(phone), phone),
  country = COALESCE(sqlc.narg(country), country),
  is_verified = COALESCE(sqlc.narg(is_verified), is_verified)
WHERE 
  id = sqlc.arg(id)
RETURNING *;