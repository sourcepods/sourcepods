package storage

import (
	"bytes"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
)

func TestInterface(t *testing.T) {
	assert.Implements(t, (*Storage)(nil), &LocalStorage{})
}

func TestLoggerOption(t *testing.T) {
	ls := &LocalStorage{}
	assert.Nil(t, ls.logger)
	LoggerOption(log.NewNopLogger())(ls)
	assert.NotNil(t, ls.logger)
}

func TestRepoPath(t *testing.T) {
	ls := &LocalStorage{root: "foo"}

	ret := ls.repoPath("foo-bar-baz")
	assert.Equal(t, "foo/fo/ob/arbaz", ret)
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
		Author: Signature{
			Name:  "First Lastname",
			Email: "first.lastname@example.com",
			Date:  time.Unix(1505935797, 0).In(time.FixedZone("", -25200)),
		},
		Committer: Signature{
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

func TestParseTreeEntry(t *testing.T) {
	te, err := parseTreeEntry("100644 blob dc2a1e6aeb5b1cf6f71666e4beb410457bdc114b	Gopkg.lock")
	assert.Nil(t, err)
	assert.Equal(t, TreeEntry{
		Mode:   "100644",
		Type:   "blob",
		Object: "dc2a1e6aeb5b1cf6f71666e4beb410457bdc114b",
		Path:   "Gopkg.lock",
	}, te)

	te, err = parseTreeEntry("040000 tree da792716f0b647e79fbfbff6c2462308791a7ea7	vendor")
	assert.Nil(t, err)
	assert.Equal(t, TreeEntry{
		Mode:   "040000",
		Type:   "tree",
		Object: "da792716f0b647e79fbfbff6c2462308791a7ea7",
		Path:   "vendor",
	}, te)
}
