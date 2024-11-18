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