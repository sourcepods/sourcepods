package session

// Lookup returns the named statement.
func Lookup(name string) string {
	return index[name]
}

var index = map[string]string{
	"clear-expired": clearExpired,
	"find-by-id":    findById,
	"save":          save,
}

var clearExpired = `
DELETE FROM sessions
WHERE expires < now();
`

var findById = `
SELECT
  sessions.id,
  sessions.expires,
  users.id       AS user_id,
  users.username AS user_username
FROM sessions
  JOIN users ON sessions.owner_id = users.id
WHERE sessions.id = $1;
`

var save = `
INSERT INTO sessions (expires, owner_id) VALUES ($1, $2)
RETURNING id;
`
