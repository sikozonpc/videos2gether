package hub

import (
	"time"

	"github.com/gorilla/websocket"
)

// Time allowed to write a message to the peer.
var writeWaitTime = 10 * time.Second

const (
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Maximum message size allowed from peer.
	maxMessageSize = 512
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

// Connection is an middleman between the websocket connection and the hub.
type Connection struct {
	WS *websocket.Conn
	// Buffered channel of outbound messages.
	Send (chan []byte)
}

type Message struct {
	Data []byte
	RoomID string
}

type SocketMessage struct {
	Action   string        `json:"action"`
	Data     VideoData     `json:"data"`
	Metadata EventMetaData `json:"metadata"`
}

type RequestSocketMessage struct {
	Action string   `json:"action"`
	Data   Playlist `json:"data"`
}

// Writes a message with a given type and payload data
func (c *Connection) Write(mt int, payload []byte) error {
	dt := time.Now().Add(writeWaitTime)
	c.WS.SetWriteDeadline(dt)

	return c.WS.WriteMessage(mt, payload)
}
