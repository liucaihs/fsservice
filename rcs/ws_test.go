package main

import (
	"github.com/gorilla/websocket"
	"net/url"
	"testing"
)

var (
	u = url.URL{Scheme: "ws", Host: "localhost:9090", Path: "/rc"}
)

func Test_ws(t *testing.T) {
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatal("Dial error:", err)
	}
	defer c.Close()
}
