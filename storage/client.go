package storage

import (
	"context"

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

func (c *Client) Description(ctx context.Context, owner string, name string, description string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage.Client.Description")
	span.SetTag("owner", owner)
	span.SetTag("name", name)
	span.SetTag("description", description)
	defer span.Finish()

	_, err := c.client.Descriptions(ctx, &DescriptionRequest{
		Owner:       owner,
		Name:        name,
		Description: description,
	})
	return err
}
