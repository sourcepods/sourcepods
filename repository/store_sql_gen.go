package repository

// Lookup returns the named statement.
func Lookup(name string) string {
	return index[name]
}

var index = map[string]string{
	"create":                 create,
	"find-by-owner-and-name": findByOwnerAndName,
	"list-by-owner-id":       listByOwnerId,
	"user-id-by-username":    userIdByUsername,
}

var create = `
INSERT INTO repositories (owner_id, name, description, website, default_branch, private, bare)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, created_at, updated_at;
`

var findByOwnerAndName = `
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
`

var listByOwnerId = `
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
`

var userIdByUsername = `
SELECT id FROM users WHERE username = $1;
`
