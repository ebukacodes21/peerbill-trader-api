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