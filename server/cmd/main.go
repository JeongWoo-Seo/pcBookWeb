package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/JeongWoo-Seo/pcBookWeb/server/internal/network/http"
	"github.com/JeongWoo-Seo/pcBookWeb/server/internal/service"
)

func main() {
	port := flag.Int("port", 0, "server port")
	flag.Parse()

	service := service.NewService()
	httpServer, err := http.NewHttpNetwork(service)
	if err != nil {
		log.Fatal("failed to htttp server")
	}

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	httpServer.Start(address)
}
