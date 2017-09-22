-- name: find-by-owner-and-name

SELECT
  id,
  description,
  website,
  default_branch,
  private,
  bare,
  created_at,
  updated_at,
  owner_id
FROM repositories
WHERE
  name = $2 AND
  owner_id = (SELECT id FROM users WHERE username = $1);
