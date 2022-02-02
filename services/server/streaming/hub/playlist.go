package hub

type Playlist []VideoData

func (p Playlist) Dequeue() Playlist {
	if len(p) <= 0 {
		return p
	}

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
