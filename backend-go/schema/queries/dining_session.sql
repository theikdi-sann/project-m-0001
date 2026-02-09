-- name: CreateDiningSession :one
INSERT INTO dining_sessions (
    table_id,
    session_type_id,
    created_by,
    guest_count,
    status,
    start_time,
    expires_at,
    price_per_guest,
    total_amount
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetDiningSession :one
SELECT * FROM dining_sessions
WHERE id = $1 LIMIT 1;

-- name: GetActiveSessionByTableID :one
SELECT * FROM dining_sessions
WHERE table_id = $1 AND status = 'active'
LIMIT 1;

-- name: ListActiveDiningSessions :many
SELECT * FROM dining_sessions
WHERE status = 'active'
ORDER BY start_time DESC;

-- name: UpdateDiningSessionStatus :one
UPDATE dining_sessions
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;
