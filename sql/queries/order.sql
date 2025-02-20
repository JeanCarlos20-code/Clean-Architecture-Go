-- name: Save :exec
INSERT INTO orders (id, price, tax, final_price, issue_date, type_requisition, delete_at) 
VALUES (?,?,?,?,?,?,?);

-- name: ListCategories :many
SELECT id, price, tax, final_price, issue_date, type_requisition, delete_at 
FROM orders 
WHERE delete_at IS NULL
ORDER BY 
    CASE WHEN sqlc.arg(sort) = 'asc' THEN id END ASC,
    CASE WHEN sqlc.arg(sort) = 'desc' THEN id END DESC
LIMIT 
    CASE 
        WHEN sqlc.arg(limit) IS NULL THEN 18446744073709551615 
        ELSE sqlc.arg(limit) 
    END
OFFSET COALESCE(sqlc.arg(offset), 0);

