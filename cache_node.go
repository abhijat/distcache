package distcache

import (
	"fmt"
	"github.com/abhijat/distcache/gen"
	"golang.org/x/net/context"
	"google.golang.org/grpc/peer"
	"log"
)

type cacheNode struct {
	State node.NodeState
	Role  node.NodeRole
	cache map[string]string
	logs  []node.LogEntry
}

func (c *cacheNode) Delete(ctx context.Context, request *node.CacheDelRequest) (*node.CacheDelResponse, error) {
	delete(c.cache, request.Key)

	entry := &node.LogEntry{
		ActionType: node.ModifyActionType_DEL,
		Key:        request.Key,
	}

	c.logs = append(c.logs, *entry)
	c.PushLogEntryToPeers(entry)

	return &node.CacheDelResponse{}, nil
}

func NewCacheNode() *cacheNode {
	return &cacheNode{
		State: node.NodeState_READY,
		Role:  node.NodeRole_PEER,
		cache: make(map[string]string),
		logs:  make([]node.LogEntry, 0),
	}
}

func (c *cacheNode) Set(ctx context.Context, request *node.CacheSetRequest) (*node.CacheSetResponse, error) {
	c.cache[request.Key] = request.Value

	entry := &node.LogEntry{
		ActionType: node.ModifyActionType_SET,
		Key:        request.Key,
		Value:      request.Value,
	}

	c.logs = append(c.logs, *entry)
	c.PushLogEntryToPeers(entry)

	return &node.CacheSetResponse{Success: true}, nil
}

func (c *cacheNode) EnumerateCache(empty *node.Empty, stream node.CacheNode_EnumerateCacheServer) error {
	for key, value := range c.cache {
		if err := stream.Send(&node.CacheEntry{Key: key, Value: value}); err != nil {
			return err
		}
	}

	return nil
}

func (c *cacheNode) BecomeLeader(context.Context, *node.Empty) (*node.Empty, error) {
	log.Println("received command to become leader")
	c.Role = node.NodeRole_BECOMING_LEADER
	c.ReplayLog()

	return &node.Empty{}, nil
}

func (c *cacheNode) Get(ctx context.Context, request *node.CacheGetRequest) (*node.CacheGetResponse, error) {

	// TODO handle missing key
	response := &node.CacheGetResponse{
		Value: c.cache[request.Key],
	}

	return response, nil
}

func (c *cacheNode) SendLogEntry(ctx context.Context, entry *node.LogEntry) (*node.Empty, error) {
	c.logs = append(c.logs, *entry)
	return &node.Empty{}, nil
}

func (c *cacheNode) Ping(ctx context.Context, request *node.HeartbeatRequest) (*node.HeartbeatResponse, error) {
	p, _ := peer.FromContext(ctx)
	log.Printf("received ping request from %v\n", p.Addr)

	return &node.HeartbeatResponse{
		NodeRole:  c.Role,
		NodeState: c.State,
	}, nil
}

func (c *cacheNode) ReplayLog() {

	log.Println("beginning replay of log")

	for _, logEntry := range c.logs {

		switch logEntry.ActionType {
		case node.ModifyActionType_SET:
			log.Printf("setting value %s against key %s\n", logEntry.Value, logEntry.Key)
			c.cache[logEntry.Key] = logEntry.Value
		case node.ModifyActionType_DEL:
			log.Printf("removing key %s from log\n", logEntry.Key)
			delete(c.cache, logEntry.Key)
		}
	}

	c.Role = node.NodeRole_LEADER
	c.State = node.NodeState_READY

	log.Println("end of log replay. leader role assumed")
}

func (c *cacheNode) PushLogEntryToPeers(entry *node.LogEntry) {

	log.Printf("pushing log entry to peers: %+v\n", entry)

	for i := 0; i < 0; i++ {
		peerAddress := fmt.Sprintf("localhost:808%d", i)
		client, err := NewClient(peerAddress)

		if err != nil {
			log.Printf("failed to init client to %s, err %v\n", peerAddress, err)
			continue
		}

		_, err = client.SendLogEntry(client.ctx, entry)
		if err != nil {
			log.Printf("failed to send log entry to %s, err %s\n", peerAddress, err)
		}

		client.Close()
	}
}
