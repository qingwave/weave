package docker

import (
	"github.com/docker/docker/client"
)

const (
	defaultHost = `unix:///var/run/docker.sock`
)

type Client struct {
	client.APIClient
}

func NewClient(host string) (*Client, error) {
	if host == "" {
		host = defaultHost
	}
	cli, err := client.NewClientWithOpts(client.WithHost(host), client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &Client{APIClient: cli}, nil
}
