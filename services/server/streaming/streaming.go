package streaming

import (
	"fmt"
	"streamserver/streaming/hub"
)

func (s Socket) CreateRoom(id string) (*RoomData, error) {
	err := hub.CreateRoomPlaylist(id)
	if err != nil {
		return nil, fmt.Errorf("room already exists")
	}

	return &RoomData{ID: id}, nil
}

func (s Socket) GetRoomPlaylist(roomId string) (hub.Playlist, error) {
	p, err := hub.GetRoomPlaylist(roomId)
	return p, err
}
