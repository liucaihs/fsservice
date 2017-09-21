package remote

import (
	"log"
	"net/http"

	"sync/atomic"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Hub struct {
	nextId     uint32
	clients    map[uint32]*Client
	register   chan *Client
	unregister chan *Client
	stopped    chan bool
	Srv        *remoteServer
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[uint32]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		stopped:    make(chan bool),
	}
}

func (hub *Hub) runLoop() {
	for {
		select {
		case client := <-hub.register:
			hub.clients[client.id] = client
		case client := <-hub.unregister:
			if _, ok := hub.clients[client.id]; ok {
				delete(hub.clients, client.id)
				client.underline()

			}
		case <-hub.stopped:
			return
		}
	}
}

func (hub *Hub) stop() {
	close(hub.stopped)
}

func (hub *Hub) serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("new client coming from %v ...", conn.RemoteAddr())
	newId := atomic.AddUint32(&hub.nextId, 1)
	client := &Client{
		id:          newId,
		hub:         hub,
		conn:        conn,
		send_buffer: make(chan []byte, 32),
		ip:          conn.RemoteAddr().String(),
	}
	hub.register <- client
	go client.readPump()
	go client.writePump()
}
