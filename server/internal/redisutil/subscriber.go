package redisutil

import (
	"context"
	"log"
	"time"

	"github.com/JeongWoo-Seo/pcBookWeb/server/internal/network/ws"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Redis connection fail: %v", err)
	}

	log.Println("Redis connected")
	return rdb
}

func StartSubscriber(ctx context.Context, rdb *redis.Client, hub *ws.Hub) {
	sub := rdb.Subscribe(ctx, "laptop-updates")

	go func() {
		ch := sub.Channel()
		log.Println("[Redis] subscribed to laptop-updates")

		for msg := range ch {
			log.Println("[Redis] received:", msg.Payload)
			hub.Broadcast <- msg.Payload
		}
	}()
}
