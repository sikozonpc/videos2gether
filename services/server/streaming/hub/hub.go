package hub

import (
	"streamserver/log"

	"github.com/sirupsen/logrus"
)

// Hub maintains the set of active connections and broadcasts messages to the connections
type Hub struct {
	Connections map[string]map[*Connection]bool
	// Inbound messages from the connections.
	Broadcast (chan Message)
	// Register requests from the connections.
	Register (chan User)
	// Unregister requests from connections.
	Unregister (chan User)
	// RoomsData registered rooms data
	RoomsPlaylist map[string]Playlist
}

// Instance is the global Hub instance that manages the connected subscriptions
var Instance = Hub{
	Broadcast:     make(chan Message),
	Register:      make(chan User),
	Unregister:    make(chan User),
	Connections:   make(map[string]map[*Connection]bool),
	RoomsPlaylist: make(map[string]Playlist),
}

func (h *Hub) Run() {
	log.Logger.Info("[HUB] started")

	for {
		select {
		case u := <-h.Register:
			h.connectUser(&u)
		case u := <-h.Unregister:
			connections := h.Connections[u.RoomID]
			if connections != nil {
				if _, ok := connections[u.Conn]; ok {
					delete(connections, u.Conn)
					h.disconnectUser(&u)

					if len(connections) == 0 {
						h.removeRoom(u.RoomID)
					}
				}
			}
		case msg := <-h.Broadcast:
			chanConns := h.Connections[msg.RoomID]

			for c := range chanConns {
				select {
				case c.Send <- msg.Data:
				default:
					close(c.Send)
					delete(chanConns, c)

					if len(chanConns) == 0 {
						h.removeRoom(msg.RoomID)
					}
				}
			}
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
