package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/JeongWoo-Seo/pcBookWeb/server/internal/network/http"
)

func main() {
	port := flag.Int("port", 0, "server port")

	httpServer, err := http.NewServer()
	if err != nil {
		log.Fatal("failed to ")
	}
	address := fmt.Sprintf("0.0.0.0:%d", *port)
	httpServer.Start(address)
}
