package main

import (
	"github.com/abhijat/distcache"
	"log"
)

func main() {

	client, err := distcache.NewClient("localhost:8080")
	if err != nil {
		log.Fatalf("error while initializing client: %v\n", err)
	}

	defer client.Close()

	response, _ := distcache.PingNode(client)
	log.Println(distcache.PingResponse(response))

	err = distcache.Set(client, "foo", "bar")
	if err != nil {
		log.Fatalf("error while set: %v\n", err)
	}

	s, err := distcache.Get(client, "foo")
	if err != nil {
		log.Fatalf("failed to get: %v\n", err)
	}

	log.Printf("key from server is %s\n", s)
	distcache.BecomeLeader(client)

	response, _ = distcache.PingNode(client)
	log.Println(distcache.PingResponse(response))

	distcache.Delete(client, "foo")
	distcache.BecomeLeader(client)

	s, err = distcache.Get(client, "foo")
	log.Println(s)
}
