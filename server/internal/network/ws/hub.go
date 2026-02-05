package ws

type WSMessage struct {
	Payload  []byte
	LaptopID string
}

type Hub struct {
	Register   chan *Client
	Unregister chan *Client
	Dispatch   chan WSMessage

	BroadcastClients map[*Client]bool
	ClientsByLaptop  map[string]map[*Client]bool
}

func NewHub() *Hub {
	return &Hub{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Dispatch:   make(chan WSMessage, 256),

		BroadcastClients: make(map[*Client]bool),
		ClientsByLaptop:  make(map[string]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.Register:
			if c.mode == "broadcast" {
				h.BroadcastClients[c] = true
			} else if c.mode == "single" {
				if h.ClientsByLaptop[c.laptopID] == nil {
					h.ClientsByLaptop[c.laptopID] = make(map[*Client]bool)
				}
				h.ClientsByLaptop[c.laptopID][c] = true
			}

		case c := <-h.Unregister:
			delete(h.BroadcastClients, c)

			if m, ok := h.ClientsByLaptop[c.laptopID]; ok {
				delete(m, c)
				if len(m) == 0 {
					delete(h.ClientsByLaptop, c.laptopID)
				}
			}

			close(c.send)

		case msg := <-h.Dispatch:
			for c := range h.BroadcastClients {
				select {
				case c.send <- msg.Payload:
				default:
				}
			}

			if clients, ok := h.ClientsByLaptop[msg.LaptopID]; ok {
				for c := range clients {
					select {
					case c.send <- msg.Payload:
					default:
					}
				}
			}
		}
	}
}
