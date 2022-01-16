package hub

// Playlist is a collection of videos
type Playlist []VideoData

func (p Playlist) Unqueue() Playlist {
	_, p = p[0], p[1:]
	return p
}

func (p Playlist) Enqueue(video VideoData) Playlist {
	return append(p, video)
}

// GetCurrent returns the next video in the queue
func (p Playlist) GetCurrent() VideoData {
	if len(p) <= 0 {
		return VideoData{}
	}
	return p[0]
}

// UpdateCurrent updates the current video playing
func (p Playlist) UpdateCurrent(u VideoData) {
	if len(p) > 0 {
		p[0] = u
	}
}
