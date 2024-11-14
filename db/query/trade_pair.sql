-- name: CreateTradePair :one
INSERT INTO trade_pairs (
  trader_id, base_asset, quote_asset, buy_rate, sell_rate
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;