package storage

import (
	"context"
	"time"

	grpcopentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

// Client holds the gRPC-connection to the storage-server
type Client struct {
	repos    RepositoryClient
	branches BranchClient
	commits  CommitClient
}

// NewClient returns a new Storage client.
func NewClient(storageAddr string) (*Client, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithUnaryInterceptor(grpcopentracing.UnaryClientInterceptor()))

	conn, err := grpc.DialContext(context.Background(), storageAddr, opts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		repos:    NewRepositoryClient(conn),
		branches: NewBranchClient(conn),
		commits:  NewCommitClient(conn),
	}, nil
}

// Create a repository
func (c *Client) Create(ctx context.Context, owner, name string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Client.Create")
	span.SetTag("owner", owner)
	span.SetTag("name", name)
	defer span.Finish()

	_, err := c.repos.Create(ctx, &CreateRequest{
		Owner: owner,
		Name:  name,
	})
	return err
}

// SetDescription of a repository
func (c *Client) SetDescription(ctx context.Context, owner, name, description string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Client.Description")
	span.SetTag("owner", owner)
	span.SetTag("name", name)
	span.SetTag("description", description)
	defer span.Finish()

	_, err := c.repos.SetDescriptions(ctx, &SetDescriptionRequest{
		Owner:       owner,
		Name:        name,
		Description: description,
	})
	return err
}

// Branches returns all branches of a repository
func (c *Client) Branches(ctx context.Context, owner, name string) ([]Branch, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Client.Branches")
	span.SetTag("owner", owner)
	span.SetTag("name", name)
	defer span.Finish()

	res, err := c.branches.List(ctx, &BranchesRequest{
		Owner: owner,
		Name:  name,
	})

	var branches []Branch
	for _, b := range res.Branch {
		branches = append(branches, Branch{
			Name: b.Name,
			Sha1: b.Sha1,
			Type: b.Type,
		})
	}

	return branches, err
}

// Commit returns a single commit from a given repository
func (c *Client) Commit(ctx context.Context, owner string, name string, ref string) (Commit, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Client.Commit")
	span.SetTag("owner", owner)
	span.SetTag("name", name)
	span.SetTag("ref", ref)
	defer span.Finish()

	req := &CommitRequest{
		Owner: owner,
		Name:  name,
		Ref:   ref,
	}

	res, err := c.commits.Get(ctx, req)
	if err != nil {
		return Commit{}, err
	}

	return Commit{
		Hash:    res.GetHash(),
		Tree:    res.GetTree(),
		Parent:  res.GetParent(),
		Message: res.GetMessage(),
		Author: Signature{
			Name:  res.GetAuthor(),
			Email: res.GetAuthorEmail(),
			Date:  time.Unix(res.GetAuthorDate(), 0),
		},
		Committer: Signature{
			Name:  res.GetCommitter(),
			Email: res.GetCommitterEmail(),
			Date:  time.Unix(res.GetCommitterDate(), 0),
		},
	}, nil
}

//Tree returns the files and folders at a given ref at a path in a repository
func (c *Client) Tree(ctx context.Context, owner, name, ref, path string) ([]TreeEntry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Client.Tree")
	span.SetTag("owner", owner)
	span.SetTag("name", name)
	span.SetTag("ref", ref)
	span.SetTag("path", path)
	defer span.Finish()

	req := &TreeRequest{
		Owner: owner,
		Name:  name,
		Ref:   ref,
		Path:  path,
	}

	res, err := c.client.Tree(ctx, req)
	if err != nil {
		return nil, err
	}

	var treeEntries []TreeEntry
	for _, te := range res.TreeEntries {
		treeEntries = append(treeEntries, TreeEntry{
			Mode:   te.Mode,
			Type:   te.Type,
			Object: te.Object,
			Path:   te.Path,
		})
	}

	return treeEntries, nil
}
