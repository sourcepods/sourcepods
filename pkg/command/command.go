package command

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

var wgAll sync.WaitGroup

// WaitAll waits for all Commands to finish
// TODO: This waits for all commands to run to completion
//   This is so we don't kill any stray "git merge" or other
//   writing commands, which could leave things corrupt
func WaitAll() {
	wgAll.Wait()
}

type (
	// Command defines the interface for all shellouts
	Command interface {
		Stdout() io.Reader
		Stderr() io.Reader
		Stdin() io.Writer

		// Finish closes the Span
		//  Note, run this is you never get to Wait() (e.g. on `return err`)
		Finish()
		Wait() error
	}

	command struct {
		cmd    *exec.Cmd
		stdout io.ReadCloser
		stderr io.ReadCloser
		stdin  io.WriteCloser
		span   opentracing.Span
	}

	// Option ...
	Option func(*command) error
)

// NewSimple is for when you would usually exec.Cmd.Run() something
func NewSimple(ctx context.Context, dir, name string, args ...string) (string, error) {
	buf := &bytes.Buffer{}
	cmd, err := New(ctx, dir, name, args, StdoutWriter(buf), StderrWriter(buf))
	if err != nil {
		return "", err
	}
	err = cmd.Wait()
	return buf.String(), err
}

// New creates a new Command
//  The caller is required to call Wait() or Finish() for the tracing to work
func New(ctx context.Context, dir, name string, args []string, opts ...Option) (Command, error) {
	// NOTE: This span is Finish()ed in Wait() or Finish()...
	span, ctx := opentracing.StartSpanFromContext(ctx, "command.New")
	span.SetTag("name", name)
	span.SetTag("args", fmt.Sprintf("%q", args))
	span.SetTag("dir", dir)
	cmd := &command{
		cmd: exec.CommandContext(ctx, name, args...),
	}
	cmd.cmd.Dir = dir
	// NOTE: GIT_DIR requires abolute paths, and `dir` can be relative for now...
	//cmd.cmd.Env = append(cmd.cmd.Env, fmt.Sprintf("GIT_DIR=%s", dir))

	for _, opt := range opts {
		if err := opt(cmd); err != nil {
			span.SetTag("error", true)
			span.LogKV("error", err)
			span.Finish()
			return nil, errors.Wrap(err, "StdoutPipe")
		}
	}

	wgAll.Add(1)
	cmd.span = span
	return cmd, cmd.cmd.Start()
}

func (c *command) Stdout() io.Reader {
	return c.stdout
}

func (c *command) Stderr() io.Reader {
	return c.stderr
}

func (c *command) Stdin() io.Writer {
	return c.stdin
}

// Finish ...
//  is idempotent
func (c *command) Finish() {
	if c.span == nil {
		return
	}
	c.span.Finish()
	wgAll.Done()
	c.span = nil
}

func (c *command) Wait() error {
	defer c.Finish()

	err := c.cmd.Wait()
	if err != nil {
		c.span.SetTag("error", true)
		c.span.LogKV("error", err)
	}
	return err
}
