-- name: ListTables :many
SELECT * FROM dining_tables
ORDER BY table_number ASC;

-- name: GetTable :one
SELECT * FROM dining_tables
WHERE id = $1 LIMIT 1;
