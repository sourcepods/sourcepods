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
	if len(args) != 3 {
		return errors.New("need exactly 3 arguments: repo_hash rev path")
	}
	repoHash, rev, path := args[0], args[1], args[2]

	s, err := storage.NewLocalStorage(c.GlobalString(cmd.FlagRoot))
	if err != nil {
		return err
	}

	r, err := s.GetRepository(context.Background(), repoHash)
	if err != nil {
		return err
	}

	entries, err := r.Tree(context.Background(), rev, path)
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
