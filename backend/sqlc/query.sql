-- queries.sql
-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: CreateUser :exec
INSERT INTO users (first_name, last_name, email) VALUES ($1, $2, $3) RETURNING *;
