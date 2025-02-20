-- name: Save :exec
INSERT INTO orders (id, price, tax, final_price, issue_date, type_requisition, delete_at) 
VALUES (?,?,?,?,?,?,?);

-- name: ListCategories :many
SELECT id, price, tax, final_price, issue_date, type_requisition, delete_at 
FROM orders 

