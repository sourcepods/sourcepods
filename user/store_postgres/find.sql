-- name: find-all

SELECT
  id,
  email,
  username,
  name,
  created_at,
  updated_at
FROM users
ORDER BY name ASC;

-- name: find-by-id

SELECT
  username,
  email,
  name,
  password,
  created_at,
  updated_at
FROM users
WHERE id = $1
LIMIT 1;

-- name: find-by-username

SELECT
  id,
  email,
  name,
  password,
  created_at,
  updated_at
FROM users
WHERE username = $1
LIMIT 1;

-- name: find-by-email

SELECT
  id,
  username,
  name,
  password,
  created_at,
  updated_at
FROM users
WHERE email = $1
LIMIT 1;
