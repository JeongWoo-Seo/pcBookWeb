package ws

import (
	"context"
	"log"
	"math/rand"
	"time"

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
				log.Println("client disconnected:", c.laptopID)
				// 더이상 구독 중인 front가 없을 경우 redis 구독 goroutine을 삭제함
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
						// client가 느려 send chan이 다 차면, 처리 속도가 느린 클라이언트로 인식하여 연결을 종료
						close(c.send)
						delete(clients, c)
					}
				}
			}

		case <-ctx.Done():
			h.closeAllSubs()
			return
		}
	}
}

// 첫번째 front가 구독할 때 redis 구독  채널을 생성
// 이후 추가되는 front들을 이미 생성된 채널을 통해 정보를 얻음
func (h *Hub) ensureSubscribe(ctx context.Context, laptopID string) {
	if _, exists := h.RedisSubs[laptopID]; exists {
		return
	}

	go h.subscribe(ctx, laptopID)
}

// redis 메시지 수신 goroutine
func (h *Hub) subscribe(ctx context.Context, laptopID string) {
	channel := "laptop:" + laptopID + ":metrics"
	retryTime := time.Second

	for {
		//laptop에 대해 구독 중인 클라이언트가 없다먼 종료
		if _, ok := h.LaptopClients[laptopID]; !ok {
			return
		}

		sub := h.RDB.Subscribe(ctx, channel)
		h.RedisSubs[laptopID] = sub
		log.Println("redis subcribed:", laptopID)

		err := h.consumeRedis(ctx, sub, laptopID)
		sub.Close()
		delete(h.RedisSubs, laptopID)

		if ctx.Err() != nil {
			return
		}

		if _, ok := h.LaptopClients[laptopID]; !ok {
			return
		}

		log.Println("redis disconnected:", laptopID, "err:", err)
		log.Println("retry after:", retryTime)

		time.Sleep(withJitter(retryTime)) //redis 재열결시 모든 채널이 동시에 연결 시도하지 않도록 랜덤한 시차를 두도록함

		retryTime *= 2
		if retryTime >= 30*time.Second {
			retryTime = 30 * time.Second
		}
	}
}

func (h *Hub) consumeRedis(ctx context.Context, sub *redis.PubSub, laptopID string) error {
	for {
		msg, err := sub.ReceiveMessage(ctx)
		if err != nil {
			if err == redis.ErrClosed { //sub.Close()과 같은 PubSub 객체가 닫혔다는 때
				return err
			}
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return err
		}

		select {
		case h.Dispatch <- RedisMessage{
			LaptopID: laptopID,
			Payload:  []byte(msg.Payload),
		}:

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// 구독 중인 다른 front가 있다면 Subscribers에 구독 리스트를 삭제하지 않고
// 구독이 0 이 되는 Subscribers에 대해서만 구독 리스트에서 삭제
func (h *Hub) unsubscribe(laptopID string) {
	if sub, ok := h.RedisSubs[laptopID]; ok {
		sub.Close()
		delete(h.RedisSubs, laptopID)
		log.Println("unsubscribed:", laptopID)
	}
}

func (h *Hub) closeAllSubs() {
	for laptopID, sub := range h.RedisSubs {
		sub.Close()
		delete(h.RedisSubs, laptopID)
	}
}

func withJitter(d time.Duration) time.Duration {
	n := rand.Int63n(int64(d) / 2)
	return d + time.Duration(n)
}
