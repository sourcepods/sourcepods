-- name: save

INSERT INTO sessions (expires, owner_id) VALUES ($1, $2)
RETURNING id;
