-- name: find-by-id

SELECT
  sessions.id,
  sessions.expires,
  users.id       AS user_id,
  users.username AS user_username
FROM sessions
  JOIN users ON sessions.owner_id = users.id
WHERE sessions.id = $1;
