-- name: GetSessionType :one
SELECT * FROM dining_session_types
WHERE id = $1 LIMIT 1;
