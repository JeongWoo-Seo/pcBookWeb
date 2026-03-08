package ws

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
)

type Hub struct {
	Register        chan *Client
	Unregister      chan *Client
	ClientsByLaptop map[string]map[*Client]bool
	Subscribers     map[string]*redis.PubSub
	RDB             *redis.Client

	mu sync.RWMutex
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
			h.mu.Lock()

			if h.ClientsByLaptop[c.laptopID] == nil {
				h.ClientsByLaptop[c.laptopID] = make(map[*Client]bool)
			}

			h.ClientsByLaptop[c.laptopID][c] = true

			first := len(h.ClientsByLaptop[c.laptopID]) == 1
			h.mu.Unlock()

			if first {
				h.ensureSubscribe(ctx, c.laptopID)
			}

		case c := <-h.Unregister:
			h.mu.Lock()

			if clients, ok := h.ClientsByLaptop[c.laptopID]; ok {
				delete(clients, c)

				if len(clients) == 0 {
					delete(h.ClientsByLaptop, c.laptopID)
					h.mu.Unlock()

					h.unsubscribe(c.laptopID)
				} else {
					h.mu.Unlock()
				}
			} else {
				h.mu.Unlock()
			}

			close(c.send)
		}
	}
}

// 첫번째 front가 구독할 때 redis 구독  채널을 생성
// 이후 추가되는 front들을 생성된 채널에서 웹소켓을 통해 정보를 얻음
func (h *Hub) ensureSubscribe(ctx context.Context, laptopID string) {
	h.mu.Lock()
	if _, exists := h.Subscribers[laptopID]; exists {
		h.mu.Unlock()
		return
	}

	channel := "laptop:" + laptopID + ":metrics"
	sub := h.RDB.Subscribe(ctx, channel)
	h.Subscribers[laptopID] = sub
	h.mu.Unlock()

	go h.consumeRedis(sub, laptopID)
}

func (h *Hub) consumeRedis(sub *redis.PubSub, laptopID string) {
	ch := sub.Channel()

	for msg := range ch {

		// 1️⃣ clients snapshot 복사
		h.mu.RLock()
		clientsMap := h.ClientsByLaptop[laptopID]

		clients := make([]*Client, 0, len(clientsMap))
		for c := range clientsMap {
			clients = append(clients, c)
		}
		h.mu.RUnlock()

		// 2️⃣ 락 없이 전송
		for _, c := range clients {
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
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.ClientsByLaptop[laptopID]; ok && len(clients) > 0 {
		return
	}

	if sub, ok := h.Subscribers[laptopID]; ok {
		sub.Close()
		delete(h.Subscribers, laptopID)
	}
}
