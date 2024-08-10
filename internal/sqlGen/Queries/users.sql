-- name: CreateUser :exec
INSERT INTO users (id, user_name, passwd)
VALUES ($1, $2, $3);

-- name: GetUserByUsername :one
SELECT * FROM users WHERE user_name = $1;