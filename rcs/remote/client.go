package remote

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait = 30 * time.Second

	pingWait = 60 * time.Second

	maxMessageSize = 512
)

type Client struct {
	id          uint32
	hub         *Hub
	conn        *websocket.Conn
	send_buffer chan []byte
	ip          string
}

func (c *Client) send(buffer []byte) {
	c.send_buffer <- buffer
}

func (c *Client) underline() {
	cliExt := c.hub.Srv.model.GetClintOne(c.id)
	cliExt["type"] = "offline"
	datastr, err := json.Marshal(cliExt)
	if err != nil {
		log.Println("client offline Marshal: ", err.Error())
	}
	c.hub.Srv.mq.PutMessage("onoff", "", string(datastr))
	c.hub.Srv.model.DelExists(c.id) // remove sqlite data
	close(c.send_buffer)
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)

	for {
		c.conn.SetReadDeadline(time.Now().Add(pingWait))
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %d %v", c.id, err)
			}
			break
		}

		handleMessage(c, message)
	}
	log.Printf("client %d quit readPump loop", c.id)
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
loop:
	for {
		select {
		case message, ok := <-c.send_buffer:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				break loop
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Println("NextWriter err:", err)
				break loop
			}
			w.Write(message)
			n := len(c.send_buffer)
			for i := 0; i < n; i++ {
				w.Write(<-c.send_buffer)
			}
			if err := w.Close(); err != nil {
				log.Println("Writer Close err:", err)
				break loop
			}
		}
	}
	log.Printf("client %d quit writePump loop", c.id)
}
