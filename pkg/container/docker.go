package container

import (
	"github.com/docker/docker/client"
)

const (
	defaultHost = `unix:///var/run/docker.sock`
)

type Client struct {
	client.APIClient
}

func NewClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(client.WithHost(defaultHost), client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &Client{APIClient: cli}, nil
}
