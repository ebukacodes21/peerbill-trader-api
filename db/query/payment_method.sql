-- name: CreatePaymentMethod :one
INSERT INTO payment_methods (
  username, trade_pair_id, bank_name, account_holder, account_number, wallet_address, crypto, fiat
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetPaymentMethods :many
SELECT * FROM payment_methods 
WHERE username = $1
ORDER BY id;

-- name: UpdatePaymentMethod :exec
UPDATE payment_methods
SET
  crypto = COALESCE(sqlc.narg(crypto), crypto),
  fiat = COALESCE(sqlc.narg(fiat), fiat),
  account_number = COALESCE(sqlc.narg(account_number), account_number),
  account_holder = COALESCE(sqlc.narg(account_holder), account_holder),
  bank_name = COALESCE(sqlc.narg(bank_name), bank_name),
  wallet_address = COALESCE(sqlc.narg(wallet_address), wallet_address)
WHERE 
  id = sqlc.arg(id)
AND username = sqlc.arg(username);

-- name: DeletePaymentMethod :exec
DELETE FROM payment_methods
WHERE id = $1
AND username = $2;