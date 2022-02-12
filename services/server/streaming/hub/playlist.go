package hub

import (
	"encoding/json"
)

type Playlist []VideoData

func unmarshalPlaylist(marshalled interface{}) Playlist {
	res := Playlist{}
	bytes, _ := marshalled.([]byte)

	err := json.Unmarshal(bytes, &res)
	if err != nil {
		panic(err)
	}

	return res
}
