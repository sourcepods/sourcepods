package mux_test

import (
	"context"
	"testing"

	"github.com/go-kit/kit/log/level"

	"github.com/stretchr/testify/assert"

	"github.com/gliderlabs/ssh"
	"github.com/sourcepods/sourcepods/pkg/ssh/mux"
)

type NopLogger struct {
	kvs []interface{}
}

func (n *NopLogger) Log(kvs ...interface{}) error {
	n.kvs = append(n.kvs, kvs...)
	return nil
}

func TestRecoverWare(t *testing.T) {
	nLogger := &NopLogger{}
	rw := mux.RecoverWare(nLogger)

	noPanic := mux.HandlerFunc(func(context.Context, ssh.Session) error {
		return nil
	})
	willPanic := mux.HandlerFunc(func(context.Context, ssh.Session) error {
		panic("foo")
	})

	assert.NotPanics(t, func() { noPanic(nil, nil) })
	assert.PanicsWithValue(t, "foo", func() { willPanic(nil, nil) })
	assert.NotPanics(t, func() { rw(nil, noPanic, nil) })
	assert.NotPanics(t, func() { rw(nil, willPanic, nil) })

	assert.Len(t, nLogger.kvs, 6)
	assert.EqualValues(t, level.ErrorValue(), nLogger.kvs[1])
	assert.EqualValues(t, "handler paniced", nLogger.kvs[3])
	assert.EqualValues(t, "foo", nLogger.kvs[5])
}
