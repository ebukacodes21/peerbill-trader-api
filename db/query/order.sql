-- name: CreateOrder :one
INSERT INTO orders (
  username, escrow_address, user_address, order_type, crypto, fiat, crypto_amount, fiat_amount, rate, bank_name, account_number, account_holder, duration
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
)
RETURNING *;

-- name: GetOrders :many
SELECT * FROM orders
WHERE username = $1
ORDER BY id;

-- name: GetUserOrders :many
SELECT * FROM orders
WHERE user_address = $1
ORDER BY id;

-- name: GetOrder :one
SELECT * FROM orders
WHERE id = $1
AND order_type = $2
LIMIT 1;

-- name: UpdateOrder :exec
UPDATE orders
SET
  escrow_address = COALESCE(sqlc.narg(escrow_address), escrow_address),
  user_address = COALESCE(sqlc.narg(user_address), user_address),
  order_type = COALESCE(sqlc.narg(order_type), order_type),
  crypto = COALESCE(sqlc.narg(crypto), crypto),
  fiat = COALESCE(sqlc.narg(fiat), fiat),
  crypto_amount = COALESCE(sqlc.narg(crypto_amount), crypto_amount),
  fiat_amount = COALESCE(sqlc.narg(fiat_amount), fiat_amount),
  rate = COALESCE(sqlc.narg(rate), rate),
  is_accepted = COALESCE(sqlc.narg(is_accepted), is_accepted),
  is_completed = COALESCE(sqlc.narg(is_completed), is_completed),
  is_rejected = COALESCE(sqlc.narg(is_rejected), is_rejected),
  is_received = COALESCE(sqlc.narg(is_received), is_received),
  duration = COALESCE(sqlc.narg(duration), duration)
WHERE 
  id = sqlc.arg(id)
  AND username = sqlc.arg (username);