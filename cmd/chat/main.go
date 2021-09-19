package main

import (
	"fmt"
	"log"
	"net/http"

	httpChat "github.com/tyghr/social_network/internal/httpserver/chat"
	redisChat "github.com/tyghr/social_network/internal/storage/chat"
)

var (
	listenPort    = 8080
	redisNodes    = []string{"redis_node_0:6379", "redis_node_1:6379", "redis_node_2:6379", "redis_node_3:6379", "redis_node_4:6379", "redis_node_5:6379"}
	redisPassword = "bitnami"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	chat := httpChat.Init(
		redisChat.Init(redisNodes, redisPassword),
	)
	log.Println("start listening...")
	if err := http.ListenAndServe(
		fmt.Sprintf("0.0.0.0:%d", listenPort),
		chat,
	); err != nil {
		log.Fatal(err)
	}
}
