-- name: CreateTransfer :one
INSERT INTO transfers (
  to_account_id, from_account_id, amount
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = $1 LIMIT 1;

-- name: ListTransfers :many
SELECT * FROM transfers
ORDER BY id
Limit $1
OFFSET $2
;
