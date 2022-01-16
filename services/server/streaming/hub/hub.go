package hub

import (
	"streamserver/log"

	"github.com/sirupsen/logrus"
)

// Hub maintains the set of active connections and broadcasts messages to the connections
type Hub struct {
	// Registered connections.
	Rooms map[string]map[*Connection]bool
	// Inbound messages from the connections.
	Broadcast chan Message
	// Register requests from the connections.
	Register chan Subscription
	// Unregister requests from connections.
	Unregister chan Subscription
	// RoomsData registred rooms data
	RoomsPlaylist map[string]Playlist
}

// Instance is the global Hub instance that manages the connected subscriptions
var Instance = Hub{
	Broadcast:     make(chan Message),
	Register:      make(chan Subscription),
	Unregister:    make(chan Subscription),
	Rooms:         make(map[string]map[*Connection]bool),
	RoomsPlaylist: make(map[string]Playlist),
}

// Run the Hub instance
func (h *Hub) Run() {
	log.Logger.Info("[HUB] started")

	for {
		select {
		case sub := <-h.Register:
			connections := h.Rooms[sub.Room]
			if connections == nil {
				connections = make(map[*Connection]bool)
				h.Rooms[sub.Room] = connections
			}
			h.Rooms[sub.Room][sub.Conn] = true

			log.Logger.WithFields(logrus.Fields{
				"room":     sub.Room,
				"room_len": len(h.Rooms[sub.Room]),
			}).Info("[HUB] Client joined room")

		case sub := <-h.Unregister:
			log.Logger.WithFields(logrus.Fields{
				"room":     sub.Room,
				"room_len": len(h.Rooms[sub.Room]),
			}).Info("[HUB] Client left the room")

			connections := h.Rooms[sub.Room]
			if connections != nil {
				if _, ok := connections[sub.Conn]; ok {
					delete(connections, sub.Conn)

					// No more users in the room
					if len(h.Rooms[sub.Room]) <= 0 {
						h.deleteRoom(sub)
					}

					close(sub.Conn.Send)
					if len(connections) == 0 {
						delete(h.Rooms, sub.Room)
					}
				}
			}
		case m := <-h.Broadcast:
			connections := h.Rooms[m.Room]

			for c := range connections {
				select {
				case c.Send <- m.Data:
				default:
					close(c.Send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(h.Rooms, m.Room)
					}
				}
			}
		}
	}
}

// CheckRoomAvailability checks if a room exists
func CheckRoomAvailability(id string) bool {
	connections := Instance.Rooms[id]
	return len(connections) > 0
}

func (h *Hub) deleteRoom(s Subscription) {
	log.Logger.WithFields(logrus.Fields{
		"room":     s.Room,
		"room_len": len(h.Rooms[s.Room]),
	}).Info("[HUB] Room deleted")

	delete(h.RoomsPlaylist, s.Room)
}
