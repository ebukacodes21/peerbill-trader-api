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
ORDER BY id
LIMIT $3
OFFSET $4;