package ws

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		mode := c.DefaultQuery("mode", "broadcast")
		laptopID := c.Query("id")

		if mode == "single" && laptopID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "id is required for single mode",
			})
			return
		}

		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("upgrade error:", err)
			return
		}

		client := &Client{
			hub:      hub,
			conn:     conn,
			send:     make(chan []byte, 32),
			laptopID: laptopID,
			mode:     mode,
		}

		hub.Register <- client
		log.Println("register sent:", client)

		go writePump(client)
		go readPump(client)
	}
}
