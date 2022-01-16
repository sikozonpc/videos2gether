package streaming

// Initialize streaming application service
func Initialize() Socket {
	return Socket{}
}

// Service represents auth service interface
type Service interface {
	CreateRoom(id string) (*RoomData, error)
	GetRoomPlaylist(roomID string) []string
}

// Socket represents streaming application service
type Socket struct{}

// RoomData struct
type RoomData struct {
	ID string
}
