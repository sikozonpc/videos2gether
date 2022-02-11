package hub

import (
	"encoding/json"
	"streamserver/log"
	_redis "streamserver/redis"

	"github.com/go-redis/redis"
)

type Playlist []VideoData

type InitialRoomPlaylist struct {
	Playlist Playlist `json:"playlist"`
}

func CheckIfRoomExists(roomId string) bool {
	_, err := _redis.Client.JSONGet(roomId, ".")
	switch {
	case err == redis.Nil:
		log.Logger.Errorf("key does not exist", err)
		return false
	case err != nil:
		log.Logger.Errorf("key does not exist", err)
		return false
	}

	return true
}

func DeleteRoom(roomId string, rdb *redis.Client) error {
	_, err := _redis.Client.JSONDel(roomId, ".playlist")
	return err
}

func GetRoomPlaylist(roomId string) (Playlist, error) {
	p, err := _redis.Client.JSONGet(roomId, ".playlist")
	if err != nil {
		return Playlist{}, err
	}

	up := unmarshalPlaylist(p)

	return up, nil
}

func CreateRoomPlaylist(roomId string) error {
	_, err := _redis.Client.JSONSet(roomId, ".", InitialRoomPlaylist{Playlist{}})
	if err != nil {
		return err
	}
	return nil
}

func GetCurrentVideo(roomId string) (VideoData, error) {
	vid, err := _redis.Client.JSONGet(roomId, ".playlist[0]")
	if err != nil {
		return VideoData{}, err
	}

	uv := unmarshalVideo(vid)

	return uv, nil
}

func AddVideoToPlaylist(roomId string, vid VideoData) error {
	_, err := _redis.Client.JSONArrAppend(roomId, ".playlist", vid)
	return err
}

func ShiftPlaylistVideo(roomId string) {
	_, err := _redis.Client.JSONArrTrim(roomId, ".playlist", 1, -1)
	if err != nil {
		log.Logger.Error(err)
		return
	}
}

func UpdateVideo(roomId string, index int, updatedVid VideoData) {
	currPl, _ := GetRoomPlaylist(roomId)
	currPl[index] = updatedVid

	_, err := _redis.Client.JSONSet(roomId, ".playlist", currPl)
	if err != nil {
		log.Logger.Error(err)
		return
	}
}

func unmarshalPlaylist(marshalled interface{}) Playlist {
	res := Playlist{}
	bytes, _ := marshalled.([]byte)

	err := json.Unmarshal(bytes, &res)
	if err != nil {
		panic(err)
	}

	return res
}
