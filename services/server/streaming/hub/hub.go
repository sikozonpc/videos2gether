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
	// RoomsData is the rooms state data
	RoomsPlaylist map[string]Playlist
	// Connections are a map of connected clients for each room
	Connections   map[string]map[*Connection]bool
}

// Instance is the global Hub instance that manages the holds the application state
var Instance = Hub{
	Broadcast:     make(chan Message),
	Connect:       make(chan User),
	Disconnect:    make(chan User),
	Connections:   make(map[string]map[*Connection]bool),
	RoomsPlaylist: make(map[string]Playlist),
}

func (h *Hub) Run() {
	log.Logger.Info("[HUB] started")

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

func CheckRoomAvailability(id string) bool {
	connections := Instance.Connections[id]
	return len(connections) > 0
}

func (h *Hub) connectUser(u *User) {
	currRoomConnections := h.Connections[u.RoomID]
	if currRoomConnections == nil {
		currRoomConnections = make(map[*Connection]bool)
		h.Connections[u.RoomID] = currRoomConnections
	}

	h.Connections[u.RoomID][u.Conn] = true

	u.syncToRoom()

	log.Logger.WithFields(logrus.Fields{
		"user":     u.Name,
		"room":     u.RoomID,
		"room_len": len(h.Connections[u.RoomID]),
	}).Info("[HUB] Client joined room")
}

func (h *Hub) disconnectUser(u *User) {
	u.disconnect()

	log.Logger.WithFields(logrus.Fields{
		"user":     u.Name,
		"room":     u.RoomID,
		"room_len": len(h.Connections[u.RoomID]),
	}).Info("[HUB] User disconnected")
}

func (h *Hub) handleDisconnectUser(u *User) {
	connections := h.Connections[u.RoomID]
	if connections != nil {
		if _, ok := connections[u.Conn]; ok {
			delete(connections, u.Conn)
			h.disconnectUser(u)

			if len(connections) == 0 {
				h.removeRoom(u.RoomID)
			}
		}
	}
}

func (h *Hub) removeRoom(roomID string) {
	delete(h.RoomsPlaylist, roomID)
	delete(h.Connections, roomID)

	log.Logger.WithFields(logrus.Fields{
		"room": roomID,
	}).Info("[HUB] Room deleted")
}

func (h *Hub) getRoomPop(id string) int {
	return len(h.Connections[id])
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
				h.removeRoom(msg.RoomID)
			}
		}
	}
}
