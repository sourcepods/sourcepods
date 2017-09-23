-- name: update
UPDATE users
SET username = $2, email = $3, name = $4, updated_at = now()
WHERE id = $1;
