package streaming

import "streamserver/streaming/hub"

// Initialize streaming application service
func Initialize() Socket {
	return Socket{}
}

// Service represents auth service interface
type Service interface {
	CreateRoom(id string) (*RoomData, error)
	GetRoomPlaylist(roomID string) hub.Playlist
}

// Socket represents streaming application service
type Socket struct{}

// RoomData struct
type RoomData struct {
	ID string
}
