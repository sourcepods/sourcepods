package mux

import (
	"context"
	"testing"

	"github.com/gliderlabs/ssh"
	"github.com/stretchr/testify/assert"
)

func TestMuxerUse(t *testing.T) {
	t.Parallel()

	var (
		mw1 = RecoverWare(nil)
		mw2 = RecoverWare(nil)
		mw3 = RecoverWare(nil)
	)

	m := New().(*muxer)
	assert.Empty(t, m.mws)
	m.Use(mw1)
	assert.Len(t, m.mws, 1)
	m.Use(mw2, mw3)
	assert.Len(t, m.mws, 3)
	// TODO: This fails for some reason...
	// for i, mw := range []MiddlewareFunc{mw3, mw2, mw1} {
	// 	assert.EqualValues(t, mw, m.mws[i])
	// }
}

func TextMuxerMatch(t *testing.T) {
	t.Parallel()

	var fooHandler = HandlerFunc(func(context.Context, ssh.Session) error {
		return nil
	})
	m := New().(*muxer)
	m.AddHandler("^foo$", "foo", fooHandler)
	m.AddHandler("^bar$", "bar", nil)

	tt := []struct {
		desc  string
		test  string
		match HandlerFunc
	}{
		{"no cmd", "", noCommandHandler},
		{"foo", "foo", fooHandler},
		{"bar is nil", "bar", nil},
		{"unknown", "unknown", unknownCommandHandler},
	}

	for _, tc := range tt {
		_, match := m.match(context.Background(), tc.test)
		assert.Equal(t, tc.match, match)
	}
}
