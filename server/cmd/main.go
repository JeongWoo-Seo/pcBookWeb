package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/JeongWoo-Seo/pcBookWeb/server/internal/network/http"
	"github.com/JeongWoo-Seo/pcBookWeb/server/internal/network/ws"
	"github.com/JeongWoo-Seo/pcBookWeb/server/internal/redisutil"
	"github.com/JeongWoo-Seo/pcBookWeb/server/internal/service"
)

func main() {
	port := flag.Int("port", 0, "server port")
	flag.Parse()

	rdb := redisutil.NewRedisClient()
	defer rdb.Close()

	service := service.NewService(rdb)
	hub := ws.NewHub(rdb)
	go hub.Run(context.Background())

	httpServer, err := http.NewHttpNetwork(service, rdb, hub)
	if err != nil {
		log.Fatal("failed to htttp server")
	}

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	httpServer.Start(address)
}
