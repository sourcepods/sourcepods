package command

import "io"

func StdinPipe(c *command) (err error) {
	c.stdin, err = c.cmd.StdinPipe()
	return
}

func StdinWriter(b io.Reader) Option {
	return func(c *command) error {
		c.cmd.Stdin = b
		return nil
	}
}

func StdoutPipe(c *command) (err error) {
	c.stdout, err = c.cmd.StdoutPipe()
}

func StdoutWriter(b io.Writer) Option {
	return func(c *command) error {
		c.cmd.Stdout = b
		return nil
	}
}

func StderrPipe(c *command) (err error) {
	c.stderr, err = c.cmd.StderrPipe()
}

func StderrWriter(b io.Writer) Option {
	return func(c *command) error {
		c.cmd.Stderr = b
		return nil
	}
}
