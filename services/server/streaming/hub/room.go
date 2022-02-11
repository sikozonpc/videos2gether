package hub

type Room struct {
	Id string
	Playlist Playlist `json:"playlist"`
}