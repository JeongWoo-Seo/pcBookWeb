package ws

import "github.com/gorilla/websocket"

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	laptopID string
	mode     string
}

func readPump(c *Client) {
	defer func() {
		c.hub.Unregister <- c
		c.conn.Close()
	}()

	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			break
		}
	}
}

func writePump(c *Client) {
	defer c.conn.Close()

	for msg := range c.send {
		c.conn.WriteMessage(websocket.TextMessage, []byte(msg))
	}
}
