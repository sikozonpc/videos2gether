package streaming

import (
	"fmt"
	"streamserver/streaming/hub"
)

func (s Socket) CreateRoom(id string) (*RoomData, error) {
	roomExists := hub.CheckRoomAvailability(id)
	if roomExists {
		return nil, fmt.Errorf("room already exists")
	}

	return &RoomData{ID: id}, nil
}

func (s Socket) GetRoomPlaylist(roomID string) hub.Playlist {
	return hub.Instance.RoomsPlaylist[roomID]
}
