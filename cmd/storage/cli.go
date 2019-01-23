package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/sourcepods/sourcepods/cmd"
	"github.com/sourcepods/sourcepods/pkg/storage"
	"github.com/urfave/cli"
)

func treeAction(c *cli.Context) error {
	args := c.Args()
	if len(args) != 4 {
		return errors.New("need exactly 4 arguments: owner name rev path")
	}
	owner, name, rev, path := args[0], args[1], args[2], args[3]

	s, err := storage.NewStorage(c.GlobalString(cmd.FlagRoot))
	if err != nil {
		return err
	}

	entries, err := s.Tree(context.Background(), owner, name, rev, path)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 3, 0, 3, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "Mode\tType\tObject\tPath")
	for _, e := range entries {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", e.Mode, e.Type, e.Object, e.Path)
	}

	return w.Flush()
}
