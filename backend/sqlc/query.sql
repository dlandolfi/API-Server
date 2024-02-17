-- queries.sql

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByLastName :one
SELECT * FROM users WHERE last_name = $1;

-- name: GetAllUsers :many
SELECT * FROM users;

-- name: CreateUser :one
INSERT INTO users (first_name, last_name, email) VALUES ($1, $2, $3) RETURNING *;
