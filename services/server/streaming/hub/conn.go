package hub

import (
	"time"

	"github.com/gorilla/websocket"
)

// Time allowed to write a message to the peer.
var writeWait = 10 * time.Second

// Connection is an middleman between the websocket connection and the hub.
type Connection struct {
	Conn *websocket.Conn
	// Buffered channel of outbound messages.
	Send chan []byte
}

// Write writes a message with a given type and payload data
func (c *Connection) Write(mt int, payload []byte) error {
	dt := time.Now().Add(writeWait)
	c.Conn.SetWriteDeadline(dt)

	return c.Conn.WriteMessage(mt, payload)
}
