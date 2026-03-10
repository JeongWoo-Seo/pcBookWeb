package ws

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisMessage struct {
	LaptopID string
	Payload  []byte
}

type Hub struct {
	RDB *redis.Client

	Register   chan *Client
	Unregister chan *Client
	Dispatch   chan RedisMessage

	LaptopClients map[string]map[*Client]bool // laptop을 구독 중인 front list
	RedisSubs     map[string]*redis.PubSub    // 구독 중인 laptop list //구독 중인 -> 정보를 받고 있는
}

func NewHub(rdb *redis.Client) *Hub {
	return &Hub{
		RDB: rdb,

		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Dispatch:   make(chan RedisMessage, 1024),

		LaptopClients: make(map[string]map[*Client]bool),
		RedisSubs:     make(map[string]*redis.PubSub),
	}
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {

		case c := <-h.Register:

			if h.LaptopClients[c.laptopID] == nil {
				h.LaptopClients[c.laptopID] = make(map[*Client]bool)
			}

			h.LaptopClients[c.laptopID][c] = true

			// 첫번째 client이면 redis subscribe
			if len(h.LaptopClients[c.laptopID]) == 1 {
				h.ensureSubscribe(ctx, c.laptopID)
			}

		case c := <-h.Unregister:

			if clients, ok := h.LaptopClients[c.laptopID]; ok {

				delete(clients, c)

				// 더이상 구독 중인 front가 없을 경우 redis 구독 채널을 삭제함
				if len(clients) == 0 {
					delete(h.LaptopClients, c.laptopID)
					h.unsubscribe(c.laptopID)
				}
			}

			close(c.send)

		//ws 메시지 보내기
		case msg := <-h.Dispatch:

			if clients, ok := h.LaptopClients[msg.LaptopID]; ok {
				for c := range clients {
					select {
					case c.send <- msg.Payload:
					default:
					}
				}
			}

		case <-ctx.Done():
			return
		}
	}
}

// 첫번째 front가 구독할 때 redis 구독  채널을 생성
// 이후 추가되는 front들을 생성된 채널에서 웹소켓을 통해 정보를 얻음
func (h *Hub) ensureSubscribe(ctx context.Context, laptopID string) {
	if _, exists := h.RedisSubs[laptopID]; exists {
		return
	}

	channel := "laptop:" + laptopID + ":metrics"
	sub := h.RDB.Subscribe(ctx, channel)

	h.RedisSubs[laptopID] = sub

	go h.consumeRedis(sub, laptopID)
}

// redis 메시지 수신 goroutine
func (h *Hub) consumeRedis(sub *redis.PubSub, laptopID string) {
	ch := sub.Channel()

	for msg := range ch {
		h.Dispatch <- RedisMessage{
			LaptopID: laptopID,
			Payload:  []byte(msg.Payload),
		}
	}
}

// 구독 중인 다른 front가 있다면 Subscribers에 구독 리스트를 삭제하지 않고
// 구독이 0 이 되는 Subscribers에 대해서만 구독 리스트에서 삭제
func (h *Hub) unsubscribe(laptopID string) {
	if sub, ok := h.RedisSubs[laptopID]; ok {

		sub.Close()

		delete(h.RedisSubs, laptopID)
	}
}
