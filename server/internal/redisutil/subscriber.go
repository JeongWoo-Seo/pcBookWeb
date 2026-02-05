package redisutil

import (
	"context"
	"encoding/json"
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

func StartRedisSubscriber(ctx context.Context, hub *ws.Hub, rdb *redis.Client) {
	sub := rdb.Subscribe(ctx, "laptop-metrics")
	ch := sub.Channel()

	for msg := range ch {
		var payload struct {
			ID string `json:"id"`
		}

		if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
			continue
		}

		hub.Dispatch <- ws.WSMessage{
			Payload:  []byte(msg.Payload),
			LaptopID: payload.ID,
		}
	}
}
