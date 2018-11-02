package main

import (
	"github.com/abhijat/distcache"
	"github.com/abhijat/distcache/gen"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {

	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	cacheNode := distcache.NewCacheNode()
	server := grpc.NewServer()

	node.RegisterCacheNodeServer(server, cacheNode)

	server.Serve(lis)

}
