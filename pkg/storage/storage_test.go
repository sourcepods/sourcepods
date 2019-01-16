package storage

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCommit(t *testing.T) {
	foo := `tree 40279100b292dd26bfda150adf1c4fd5a4e52ffe
parent ae51e9d1b987f9086cbc65e694f06759bc62e743
author First Lastname <first.lastname@example.com> 1505935797 -0700
committer Second Lastname <second.lastname@example.com> 1505935797 -0700

do something very useful to conquer the world

my
awesome

body`

	commit, err := parseCommit(bytes.NewBufferString(foo), "99cc2f794893815dfc69ab1ba3370ef3e7a9fed2")
	assert.NoError(t, err)
	assert.Equal(t, "99cc2f794893815dfc69ab1ba3370ef3e7a9fed2", commit.Hash)
	assert.Equal(t, "40279100b292dd26bfda150adf1c4fd5a4e52ffe", commit.Tree)
	assert.Equal(t, "ae51e9d1b987f9086cbc65e694f06759bc62e743", commit.Parent)
	assert.Equal(t, "First Lastname", commit.Author.Name)
	assert.Equal(t, "first.lastname@example.com", commit.Author.Email)
	assert.Equal(t, int64(1505935797), commit.Author.Date.Unix())
	assert.Equal(t, "Second Lastname", commit.Committer.Name)
	assert.Equal(t, "second.lastname@example.com", commit.Committer.Email)
	assert.Equal(t, int64(1505935797), commit.Committer.Date.Unix())
	assert.Equal(t, "do something very useful to conquer the world", commit.Message)
	assert.Equal(t, "\n\nmy\nawesome\n\nbody", commit.Body)
}
