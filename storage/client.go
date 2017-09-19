package storage

import (
	"context"

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

func (c *Client) Create(owner string, name string) error {
	_, err := c.client.Create(context.Background(), &CreateRequest{
		Owner: owner,
		Name:  name,
	})
	return err
}

func (c *Client) Description(owner string, name string, description string) error {
	_, err := c.client.Descriptions(context.Background(), &DescriptionRequest{
		Owner:       owner,
		Name:        name,
		Description: description,
	})
	return err
}
