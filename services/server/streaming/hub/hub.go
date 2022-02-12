package hub

import (
	"streamserver/log"

	"github.com/sirupsen/logrus"
)

// Hub represents a lobby and it's responsible for connecting, disconnecting and broadcasting
// the users events.
// It runs as a loop listenning and redirecting the web socket request to it's handler.
type Hub struct {
	// Inbound messages from the connections.
	Broadcast (chan Message)
	// Connect requests from the connections.
	Connect (chan User)
	// Disconnect requests from connections.
	Disconnect (chan User)
	// Connections are a map of connected clients for each room
	Connections map[string]map[*Connection]bool
}

// Instance is the global Hub instance that manages the holds the application state
var Instance = Hub{
	Broadcast:   make(chan Message),
	Connect:     make(chan User),
	Disconnect:  make(chan User),
	Connections: make(map[string]map[*Connection]bool),
}

func (h *Hub) Listen() {
	log.Logger.Info("[Hub] is listening")

	for {
		select {
		case u := <-h.Connect:
			h.connectUser(&u)
		case u := <-h.Disconnect:
			h.handleDisconnectUser(&u)
		case msg := <-h.Broadcast:
			h.broadcast(msg)
		}
	}
}

func GetRoomConnections(id string) int {
	return len(Instance.Connections[id])
}

func (h *Hub) connectUser(u *User) {
	roomExists := u.Room.checkIfExists()
	if !roomExists {
		log.Logger.WithFields(logrus.Fields{
			"user":     u.Name,
			"room":     u.Room.Id,
			"room_len": len(h.Connections[u.Room.Id]),
		}).Info("[Hub] Client tried to join room that does not exist")
		return
	}

	currRoomConnections := h.Connections[u.Room.Id]
	if currRoomConnections == nil {
		currRoomConnections = make(map[*Connection]bool)
		h.Connections[u.Room.Id] = currRoomConnections
	}

	h.Connections[u.Room.Id][u.Conn] = true

	log.Logger.WithFields(logrus.Fields{
		"user":     u.Name,
		"room":     u.Room.Id,
		"room_len": len(h.Connections[u.Room.Id]),
	}).Info("[Hub] Client has joined room")
}

func (h *Hub) disconnectUser(u *User) {
	u.disconnect()

	log.Logger.WithFields(logrus.Fields{
		"user":     u.Name,
		"room":     u.Room.Id,
		"room_len": len(h.Connections[u.Room.Id]),
	}).Info("[Hub] User has disconnected")
}

func (h *Hub) handleDisconnectUser(u *User) {
	connections := h.Connections[u.Room.Id]
	if connections != nil {
		if _, ok := connections[u.Conn]; ok {
			delete(connections, u.Conn)
			h.disconnectUser(u)

			if len(connections) == 0 {
				h.removeRoomConnections(u.Room.Id)
			}
		}
	}
}

func (h *Hub) removeRoomConnections(roomID string) {
	delete(h.Connections, roomID)

	log.Logger.WithFields(logrus.Fields{
		"room": roomID,
	}).Info("[Hub] Room connections got deleted")
}

func (h *Hub) broadcast(msg Message) {
	roomConns := h.Connections[msg.RoomID]

	for c := range roomConns {
		select {
		case c.Send <- msg.Data:
		default:
			close(c.Send)
			delete(roomConns, c)

			if len(roomConns) == 0 {
				h.removeRoomConnections(msg.RoomID)
			}
		}
	}
}
