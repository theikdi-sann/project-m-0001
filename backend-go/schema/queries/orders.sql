-- name: CreateOrder :one
INSERT INTO orders (dining_session_id, status, total_amount)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, menu_item_id, quantity, unit_price, notes)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders
WHERE id = $1 LIMIT 1;

-- name: ListOrdersBySession :many
SELECT * FROM orders
WHERE dining_session_id = $1
ORDER BY created_at DESC;

-- name: ListActiveOrders :many
SELECT * FROM orders
WHERE status IN ('pending', 'preparing', 'ready')
ORDER BY created_at ASC;

-- name: UpdateOrderStatus :one
UPDATE orders
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ListOrderItems :many
SELECT * FROM order_items
WHERE order_id = $1;
