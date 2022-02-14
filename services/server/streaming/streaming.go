package streaming

import (
	"fmt"
	"streamserver/streaming/hub"
)

func (s Socket) CreateRoom(id string) (*RoomData, error) {
	err := hub.Create(id)
	if err != nil {
		return nil, fmt.Errorf("room already exists")
	}

	return &RoomData{ID: id}, nil
}

func (s Socket) GetRoomPlaylist(roomId string) (hub.Playlist, error) {
	p, err := hub.Get(roomId)
	return p, err
}

func (s Socket) DeleteRoom(roomId string) error {
	return hub.Delete(roomId)
}

func (s Socket) CleanAllRooms() {
	hub.Flush()
}
