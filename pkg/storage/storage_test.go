package storage

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInterface(t *testing.T) {
	assert.Implements(t, (*Storage)(nil), &LocalStorage{})
}
func TestParseCommit(t *testing.T) {
	foo := `tree 40279100b292dd26bfda150adf1c4fd5a4e52ffe
parent ae51e9d1b987f9086cbc65e694f06759bc62e743
author First Lastname <first.lastname@example.com> 1505935797 -0700
committer Second Lastname <second.lastname@example.com> 1505935797 -0700
something Foobar

do something very useful to conquer the world

my
awesome

body`
	expected := Commit{
		Hash:   "99cc2f794893815dfc69ab1ba3370ef3e7a9fed2",
		Tree:   "40279100b292dd26bfda150adf1c4fd5a4e52ffe",
		Parent: "ae51e9d1b987f9086cbc65e694f06759bc62e743",
		Author: Author{
			Name:  "First Lastname",
			Email: "first.lastname@example.com",
			Date:  time.Unix(1505935797, 0).In(time.FixedZone("", -25200)),
		},
		Committer: Author{
			Name:  "Second Lastname",
			Email: "second.lastname@example.com",
			Date:  time.Unix(1505935797, 0).In(time.FixedZone("", -25200)),
		},
		Message: "do something very useful to conquer the world",
		Body:    "my\nawesome\n\nbody",
	}
	commit, err := parseCommit(bytes.NewBufferString(foo), "99cc2f794893815dfc69ab1ba3370ef3e7a9fed2")
	assert.NoError(t, err)
	assert.Equal(t, expected, commit)
}
