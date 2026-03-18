package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

type Client struct {
	conn   *grpc.ClientConn
	client collector.CollectorServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}
	
	return &Client{
		conn:   conn,
		client: collector.NewCollectorServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetRepositoryInfo(ctx context.Context, owner, repo string) (*collector.RepoResponse, error) {
	req := &collector.RepoRequest{
		Owner: owner,
		Repo:  repo,
	}
	
	return c.client.GetRepositoryInfo(ctx, req)
}
