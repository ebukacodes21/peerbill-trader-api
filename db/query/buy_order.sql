-- name: CreateBuyOrder :one
INSERT INTO buy_orders (
  username, wallet_address, crypto, fiat, crypto_amount, fiat_amount, rate, duration
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetBuyOrders :many
SELECT * FROM buy_orders
WHERE username = $1
ORDER BY id;

-- name: GetBuyOrder :one
SELECT * FROM buy_orders
WHERE id = $1
LIMIT 1;

-- name: UpdateBuyOrder :one
UPDATE buy_orders
SET
  wallet_address = COALESCE(sqlc.narg(wallet_address), wallet_address),
  crypto = COALESCE(sqlc.narg(crypto), crypto),
  fiat = COALESCE(sqlc.narg(fiat), fiat),
  crypto_amount = COALESCE(sqlc.narg(crypto_amount), crypto_amount),
  fiat_amount = COALESCE(sqlc.narg(fiat_amount), fiat_amount),
  rate = COALESCE(sqlc.narg(rate), rate),
  is_accepted = COALESCE(sqlc.narg(is_accepted), is_accepted),
  is_completed = COALESCE(sqlc.narg(is_completed), is_completed),
  is_rejected = COALESCE(sqlc.narg(is_rejected), is_rejected),
  is_expired = COALESCE(sqlc.narg(is_expired), is_expired),
  duration = COALESCE(sqlc.narg(duration), duration)
WHERE 
  id = sqlc.arg(id)
  AND username = sqlc.arg (username)
RETURNING *;