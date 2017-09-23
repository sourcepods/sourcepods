-- name: create

INSERT INTO repositories (owner_id, name, description, website, default_branch, private, bare)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, created_at, updated_at;
