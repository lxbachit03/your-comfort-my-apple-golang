-- name: GetUser :one
SELECT * FROM users WHERE user_uuid = $1;