package ws

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Hub struct {
	Register        chan *Client
	Unregister      chan *Client
	ClientsByLaptop map[string]map[*Client]bool
	Subscribers     map[string]*redis.PubSub
	RDB             *redis.Client
}

func NewHub(rdb *redis.Client) *Hub {
	return &Hub{
		Register:        make(chan *Client),
		Unregister:      make(chan *Client),
		ClientsByLaptop: make(map[string]map[*Client]bool),
		Subscribers:     make(map[string]*redis.PubSub),
		RDB:             rdb,
	}
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {

		case c := <-h.Register:
			if h.ClientsByLaptop[c.laptopID] == nil {
				h.ClientsByLaptop[c.laptopID] = make(map[*Client]bool)
			}

			h.ClientsByLaptop[c.laptopID][c] = true

			h.ensureSubscribe(ctx, c.laptopID)

		case c := <-h.Unregister:

			if clients, ok := h.ClientsByLaptop[c.laptopID]; ok {
				delete(clients, c)

				if len(clients) == 0 {
					delete(h.ClientsByLaptop, c.laptopID)
					h.unsubscribe(c.laptopID)
				}
			}

			close(c.send)
		}
	}
}

// 첫번째 front가 구독할 때 redis 구독  채널을 생성
// 이후 추가되는 front들을 생성된 채널에서 웹소켓을 통해 정보를 얻음
func (h *Hub) ensureSubscribe(ctx context.Context, laptopID string) {
	if _, exists := h.Subscribers[laptopID]; exists {
		return
	}

	channel := "laptop:" + laptopID + ":metrics"
	sub := h.RDB.Subscribe(ctx, channel)
	h.Subscribers[laptopID] = sub

	go h.consumeRedis(sub, laptopID)
}

func (h *Hub) consumeRedis(sub *redis.PubSub, laptopID string) {
	ch := sub.Channel()

	for msg := range ch {

		clients := h.ClientsByLaptop[laptopID]

		for c := range clients {
			select {
			case c.send <- []byte(msg.Payload):
			default:
			}
		}
	}
}

// 구독 중인 다른 front가 있다면 Subscribers에 구독 리스트를 삭제하지 않고
// 구독이 0 이 되는 Subscribers에 대해서만 구독 리스트에서 삭제
func (h *Hub) unsubscribe(laptopID string) {
	if clients, ok := h.ClientsByLaptop[laptopID]; ok && len(clients) > 0 {
		return
	}

	if sub, ok := h.Subscribers[laptopID]; ok {
		sub.Close()
		delete(h.Subscribers, laptopID)
	}
}
