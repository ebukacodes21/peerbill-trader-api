-- name: CreateTradePair :one
INSERT INTO
    trade_pairs (
        username,
        crypto,
        fiat,
        buy_rate,
        sell_rate
    )
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetTradePairs :many
SELECT * FROM trade_pairs 
WHERE crypto = $1
AND fiat = $2
ORDER BY id;

-- name: GetTraderPairs :many
SELECT * FROM trade_pairs 
WHERE username = $1
ORDER BY id;

-- name: UpdateTradePair :exec
UPDATE trade_pairs
SET
  crypto = COALESCE(sqlc.narg(crypto), crypto),
  fiat = COALESCE(sqlc.narg(fiat), fiat),
  buy_rate = COALESCE(sqlc.narg(buy_rate), buy_rate),
  sell_rate = COALESCE(sqlc.narg(sell_rate), sell_rate)
WHERE 
  id = sqlc.arg(id)
AND username = sqlc.arg(username);

-- name: DeleteTradePair :exec
DELETE FROM trade_pairs
WHERE id = $1
AND username = $2;

-- name: GetTradePair :one
SELECT * FROM trade_pairs 
WHERE crypto = $1
AND fiat = $2
AND username = $3
LIMIT 1;