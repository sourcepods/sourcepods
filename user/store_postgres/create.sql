-- name: create

INSERT INTO users (email, username, name, password) VALUES ($1, $2, $3, $4)
RETURNING id, created_at, updated_at;
