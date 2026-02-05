package ws

type WSMessage struct {
	Payload  []byte
	LaptopID string
}

type Hub struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Dispatch   chan WSMessage
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Dispatch:   make(chan WSMessage, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {

		case c := <-h.Register:
			h.Clients[c] = true

		case c := <-h.Unregister:
			if _, ok := h.Clients[c]; ok {
				delete(h.Clients, c)
				close(c.send)
			}

		case msg := <-h.Dispatch:
			for c := range h.Clients {
				if c.mode == "broadcast" || (c.mode == "single" && c.laptopID == msg.LaptopID) {
					c.send <- msg.Payload
				}
			}
		}
	}
}
