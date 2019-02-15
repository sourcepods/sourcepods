package command

import (
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
func WaitAll() {
	wgAll.Wait()
}

// Command defines the interface for all shellouts
type Command interface {
	Stdout() io.Reader
	Stderr() io.Reader
	Stdin() io.Writer

	Wait() error
}

type command struct {
	cmd    *exec.Cmd
	stdout io.ReadCloser
	stderr io.ReadCloser
	stdin  io.WriteCloser
	span   opentracing.Span
}

// New creates a new Command
func New(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer, dir, name string, args ...string) (Command, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "command.New")
	span.SetTag("name", name)
	span.SetTag("args", fmt.Sprintf("%v", args))
	span.SetTag("dir", dir)
	cmd := &command{
		cmd:  exec.CommandContext(ctx, name, args...),
		span: span,
	}
	cmd.cmd.Dir = dir
	// NOTE: GIT_DIR requires abolute paths, and `dir` is relative for now...
	//cmd.cmd.Env = append(cmd.cmd.Env, fmt.Sprintf("GIT_DIR=%s", dir))

	if stdout == nil {
		var err error
		cmd.stdout, err = cmd.cmd.StdoutPipe()
		if err != nil {
			return nil, errors.Wrap(err, "StdoutPipe")
		}
	} else {
		cmd.cmd.Stdout = stdout
	}
	if stderr == nil {
		var err error
		cmd.stderr, err = cmd.cmd.StderrPipe()
		if err != nil {
			return nil, errors.Wrap(err, "StderrPipe")
		}
	} else {
		cmd.cmd.Stderr = stderr
	}
	if stdin == nil {
		var err error
		cmd.stdin, err = cmd.cmd.StdinPipe()
		if err != nil {
			return nil, errors.Wrap(err, "StdinPipe")
		}
	} else {
		cmd.cmd.Stdin = stdin
	}

	wgAll.Add(1)
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

func (c *command) Wait() error {
	defer c.span.Finish()
	defer wgAll.Done()

	err := c.cmd.Wait()
	if err != nil {
		c.span.SetTag("error", true)
		c.span.LogKV("error", err)
	}
	return err
}
