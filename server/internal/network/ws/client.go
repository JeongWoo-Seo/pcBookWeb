package ws

import (
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	laptopID string
	mode     string
}

const pongWait = 60 * time.Second
const pingPeriod = 30 * time.Second
const writeWait = 10 * time.Second

func (c *Client) readPump() {
	// front 종료시 클라이언트 리스트에서 등록 삭제
	defer func() {
		c.hub.Unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512) //메시지 크기 제한(bytes)

	c.conn.SetReadDeadline(time.Now().Add(pongWait)) //일정 시간까지 read 가 없으면 connection 종료

	c.conn.SetPongHandler(func(string) error { //read timeout 연장
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {

		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok { //클라이언트 연결 종료 메시지
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			//WebSocket frame writer 생성 - 여러 메시지를 하나의 fram으로 전송
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message) //첫번째 메시지를 frame으로 전달

			// batching 시작
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
