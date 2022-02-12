package hub

import (
	"streamserver/log"
	_redis "streamserver/redis"

	"github.com/go-redis/redis"
)

type Room struct {
	Id string
}

type InitialRoomPlaylist struct {
	Playlist Playlist `json:"playlist"`
}

func Delete(roomId string) error {
	_, err := _redis.Client.JSONDel(roomId, ".playlist")
	return err
}

func Create(roomId string) error {
	_, err := _redis.Client.JSONSet(roomId, ".", InitialRoomPlaylist{Playlist{}})
	if err != nil {
		return err
	}
	return nil
}

func Get(roomId string) (Playlist, error) {
	p, err := _redis.Client.JSONGet(roomId, ".playlist")
	if err != nil {
		return Playlist{}, err
	}

	up := unmarshalPlaylist(p)
	return up, nil
}

func (r *Room) checkIfExists() bool {
	_, err := _redis.Client.JSONGet(r.Id, ".")
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

func (r *Room) getCurrentVideo() (VideoData, error) {
	vid, err := _redis.Client.JSONGet(r.Id, ".playlist[0]")
	if err != nil {
		return VideoData{}, err
	}

	uv := unmarshalVideo(vid)

	return uv, nil
}

func (r *Room) getPlaylist() (Playlist, error) {
	p, err := _redis.Client.JSONGet(r.Id, ".playlist")
	if err != nil {
		return Playlist{}, err
	}

	up := unmarshalPlaylist(p)

	return up, nil
}

func (r *Room) addVideoToPlaylist(vid VideoData) error {
	_, err := _redis.Client.JSONArrAppend(r.Id, ".playlist", vid)
	return err
}

func (r *Room) updateVideo(index int, updatedVid VideoData) {
	currPl, _ := r.getPlaylist()
	currPl[index] = updatedVid

	_, err := _redis.Client.JSONSet(r.Id, ".playlist", currPl)
	if err != nil {
		log.Logger.Error(err)
		return
	}
}
func (r *Room) shiftPlaylistVideo() {
	_, err := _redis.Client.JSONArrTrim(r.Id, ".playlist", 1, -1)
	if err != nil {
		log.Logger.Error(err)
		return
	}
}
