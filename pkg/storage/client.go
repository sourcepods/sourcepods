package storage

import (
	"context"
	"time"

	grpcopentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	opentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

//go:generate protoc storage.proto --go_out=plugins=grpc:.

type Client struct {
	client StorageClient
}

func NewClient(storageAddr string) (*Client, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithUnaryInterceptor(grpcopentracing.UnaryClientInterceptor()))

	conn, err := grpc.DialContext(context.Background(), storageAddr, opts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		client: NewStorageClient(conn),
	}, nil
}

func (c *Client) Create(ctx context.Context, owner string, name string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Client.Create")
	span.SetTag("owner", owner)
	span.SetTag("name", name)
	defer span.Finish()

	_, err := c.client.Create(ctx, &CreateRequest{
		Owner: owner,
		Name:  name,
	})
	return err
}

func (c *Client) SetDescription(ctx context.Context, owner string, name string, description string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Client.Description")
	span.SetTag("owner", owner)
	span.SetTag("name", name)
	span.SetTag("description", description)
	defer span.Finish()

	_, err := c.client.SetDescriptions(ctx, &SetDescriptionRequest{
		Owner:       owner,
		Name:        name,
		Description: description,
	})
	return err
}

func (c *Client) Branches(ctx context.Context, owner string, name string) ([]Branch, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Client.Branches")
	span.SetTag("owner", owner)
	span.SetTag("name", name)
	defer span.Finish()

	res, err := c.client.Branches(ctx, &BranchesRequest{
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

func (c *Client) Commit(ctx context.Context, owner string, name string, rev string) (Commit, error) {
	req := &CommitRequest{
		Owner: owner,
		Name:  name,
		Rev:   rev,
	}

	res, err := c.client.Commit(ctx, req)
	if err != nil {
		return Commit{}, err
	}

	return Commit{
		Hash:           res.GetHash(),
		Tree:           res.GetTree(),
		Parent:         res.GetParent(),
		Message:        res.GetMessage(),
		Author:         res.GetAuthor(),
		AuthorEmail:    res.GetAuthorEmail(),
		AuthorDate:     time.Unix(res.GetAuthorDate(), 0),
		Committer:      res.GetCommitter(),
		CommitterEmail: res.GetCommitterEmail(),
		CommitterDate:  time.Unix(res.GetCommitterDate(), 0),
	}, nil
}
