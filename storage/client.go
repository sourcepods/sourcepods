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

func (c *Client) Tree(ctx context.Context, owner, name, branch string) ([]TreeObject, error) {
	req := &TreeRequest{
		Owner:  owner,
		Name:   name,
		Branch: branch,
	}

	res, err := c.client.Tree(ctx, req)
	if err != nil {
		return nil, err
	}

	var objects []TreeObject
	for _, object := range res.Objects {
		objects = append(objects, TreeObject{
			Mode:   object.GetMode(),
			Type:   object.GetType(),
			Object: object.GetObject(),
			File:   object.GetFile(),
			Commit: Commit{
				Hash:           object.GetCommit().GetHash(),
				Tree:           object.GetCommit().GetTree(),
				Parent:         object.GetCommit().GetParent(),
				Subject:        object.GetCommit().GetSubject(),
				Author:         object.GetCommit().GetAuthor(),
				AuthorEmail:    object.GetCommit().GetAuthorEmail(),
				AuthorDate:     time.Unix(object.GetCommit().GetAuthorDate(), 0),
				Committer:      object.GetCommit().GetCommitter(),
				CommitterEmail: object.GetCommit().GetCommitterEmail(),
				CommitterDate:  time.Unix(object.GetCommit().GetCommitterDate(), 0),
			},
		})
	}

	return objects, nil
}
