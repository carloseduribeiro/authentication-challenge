-- name: InsertDebt :exec
INSERT INTO customer.debts (id, document, dueDate, amount, created_at)
VALUES ($1, $2, $3, $4, $5);

-- name: GetDebtsByDocument :many
SELECT *
FROM customer.debts
WHERE document = $1;