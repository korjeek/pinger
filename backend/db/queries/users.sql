-- name: CreateUser :exec
INSERT INTO users (id, email, password_hash)
VALUES ($1, $2, $3);

-- name: GetUSerByID :one
SELECT id, email, password_hash
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, email, password_hash
FROM users
WHERE email = $1;

-- name: UpdateUser :exec
UPDATE users
SET email = $2,
    password_hash = $3
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

