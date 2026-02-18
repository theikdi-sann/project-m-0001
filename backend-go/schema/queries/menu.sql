-- name: CreateMenuCategory :one
INSERT INTO menu_categories (name, sort_order)
VALUES ($1, $2)
RETURNING *;

-- name: ListMenuCategories :many
SELECT * FROM menu_categories
ORDER BY sort_order ASC;

-- name: CreateMenuItem :one
INSERT INTO menu_items (category_id, name, description, price, image_url, is_available)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListMenuItemsByCategory :many
SELECT * FROM menu_items
WHERE category_id = $1 AND is_available = TRUE
ORDER BY name ASC;

-- name: GetMenuItem :one
SELECT * FROM menu_items
WHERE id = $1 LIMIT 1;

-- name: ListAllMenuItemsByCategory :many
SELECT * FROM menu_items
WHERE category_id = $1
ORDER BY name ASC;

-- name: UpdateMenuItemAvailability :one
UPDATE menu_items
SET is_available = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;
