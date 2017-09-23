package user

// Lookup returns the named statement.
func Lookup(name string) string {
	return index[name]
}

var index = map[string]string{
	"create":           create,
	"find-all":         findAll,
	"find-by-id":       findById,
	"find-by-username": findByUsername,
	"find-by-email":    findByEmail,
	"update":           update,
}

var create = `
INSERT INTO users (email, username, name, password) VALUES ($1, $2, $3, $4)
RETURNING id, created_at, updated_at;
`

var findAll = `
SELECT
  id,
  email,
  username,
  name,
  created_at,
  updated_at
FROM users
ORDER BY name ASC;
`

var findById = `
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
`

var findByUsername = `
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
`

var findByEmail = `
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
`

var update = `
UPDATE users
SET username = $2, email = $3, name = $4, updated_at = now()
WHERE id = $1;
`
