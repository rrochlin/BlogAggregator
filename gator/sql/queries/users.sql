-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
	    $1,
	    $2,
	    $3,
	    $4
	)
	RETURNING *;

-- name: GetUserByUID :one

SELECT *
FROM users
WHERE id=$1;

-- name: GetUser :one

SELECT *
FROM users
WHERE name=$1;

-- name: TruncateTable :exec
TRUNCATE TABLE users CASCADE;

-- name: GetUsers :many
SELECT * FROM users;

