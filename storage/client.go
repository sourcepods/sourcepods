package storage

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

//go:generate protoc storage.proto --go_out=plugins=grpc:.

type Client struct {
	client StorageClient
}

func NewClient(conn *grpc.ClientConn) *Client {
	return &Client{
		client: NewStorageClient(conn),
	}
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
