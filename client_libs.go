package distcache

import (
	"context"
	"fmt"
	"github.com/abhijat/distcache/gen"
	"google.golang.org/grpc"
	"log"
	"time"
)

type client struct {
	node.CacheNodeClient
	conn    *grpc.ClientConn
	ctx     context.Context
	address string
	cancel  context.CancelFunc
}

func (c *client) Close() {
	defer c.conn.Close()
	defer c.cancel()
}

func NewClient(address string) (*client, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	c := &client{
		conn:    conn,
		ctx:     ctx,
		address: address,
		cancel:  cancel,
	}
	c.CacheNodeClient = node.NewCacheNodeClient(conn)

	return c, nil
}

func PingNode(c *client) (*node.HeartbeatResponse, error) {
	return c.Ping(c.ctx, &node.HeartbeatRequest{})
}

func PingResponse(response *node.HeartbeatResponse) string {
	return fmt.Sprintf(
		"node state: %s | node role: %s",
		response.NodeState.String(),
		response.NodeRole.String(),
	)
}

func Set(c *client, key, value string) error {
	response, err := c.Set(c.ctx, &node.CacheSetRequest{
		Key:   key,
		Value: value,
	})

	log.Printf("cache set success [%t] for %s -> %s\n", response.Success, key, value)
	return err
}

func Get(c *client, key string) (string, error) {
	response, err := c.Get(c.ctx, &node.CacheGetRequest{Key: key})
	if err != nil {
		return "", err
	}

	return response.Value, nil
}

func BecomeLeader(c *client) error {
	_, err := c.BecomeLeader(c.ctx, &node.Empty{})
	return err
}

func Delete(c *client, key string) {
	c.Delete(c.ctx, &node.CacheDelRequest{Key: key})
}
