-- name: list-by-owner-id

SELECT
  id,
  name,
  description,
  website,
  default_branch,
  private,
  bare,
  created_at,
  updated_at,
  (SELECT 42) AS stars,
  (SELECT 23) AS forks
FROM repositories
WHERE owner_id = $1
ORDER BY updated_at DESC;

-- name: user-id-by-username

SELECT id FROM users WHERE username = $1;
