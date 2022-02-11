package hub

import (
	"encoding/json"

	"streamserver/log"

	"github.com/sirupsen/logrus"
)

type EventType string

const (
	REQUEST             EventType = "REQUEST"
	PLAY_VIDEO          EventType = "PLAY_VIDEO"
	PAUSE_VIDEO         EventType = "PAUSE_VIDEO"
	SKIP_VIDEO          EventType = "SKIP_VIDEO"
	END_VIDEO           EventType = "END_VIDEO"
	SYNC                EventType = "SYNC"
	USER_JOINED         EventType = "USER_JOINED"
	USER_DISCONNECTED   EventType = "USER_DISCONNECTED"
	REQUEST_TIME        EventType = "REQUEST_TIME"
	SEND_TIME_TO_SERVER EventType = "SEND_TIME_TO_SERVER"
)

type EventMetaData struct {
	ActionFrom     string `json:"actionFrom"`
	UsersConnected int    `json:"usersConnected"`
	RoomID         string `json:"roomID"`
}

func HandleActionEvent(rawMsg []byte, u *User) {
	logger := log.Logger.WithFields(logrus.Fields{
		"room": u.RoomID,
		"user": u.Name,
	})

	p, err := GetRoomPlaylist(u.RoomID)
	if err != nil {
		return
	}

	eventType, data := unmarshalSocketMessage(rawMsg)
	hasVideos := len(p) > 0

	meta := EventMetaData{u.Name, GetRoomPop(u.RoomID), u.RoomID}

	switch eventType {
	case REQUEST:
		// TODO: Proper validate if it's a valid youtube video

		newVid := VideoData{
			Time:    0,
			Playing: true,
			Url:     data.Url,
		}

		err := AddVideoToPlaylist(u.RoomID, newVid)
		if err != nil {
			panic(err)
		}

		// TODO: Validate if END_VIDEO should be sent here?
		u.broadcastMessage(SocketMessage{"END_VIDEO", newVid, meta})

	case END_VIDEO:
		ShiftPlaylistVideo(u.RoomID)
		up, err := GetRoomPlaylist(u.RoomID)
		if err != nil {
			log.Logger.Fatal(err)
		}

		if len(up) == 0 {
			u.broadcastMessage(SocketMessage{"END_VIDEO", VideoData{}, meta})
		} else {
			u.broadcastMessage(SocketMessage{"END_VIDEO", up[0], meta})
		}
	case PLAY_VIDEO:
		if hasVideos {
			currVid := p[0]

			UpdateVideo(u.RoomID, 0, VideoData{currVid.Url, data.Time, true})

			up, err := GetRoomPlaylist(u.RoomID)
			if err != nil {
				log.Logger.Fatal(err)
			}

			if len(up) == 0 {
				u.broadcastMessage(SocketMessage{"PLAY_VIDEO", VideoData{}, meta})
			} else {
				u.broadcastMessage(SocketMessage{"PLAY_VIDEO", up[0], meta})
			}
		}
	case USER_JOINED:
		currVid := p[0]
		u.broadcastMessage(SocketMessage{"USER_JOINED", currVid, meta})
	case USER_DISCONNECTED:
		currVid := p[0]
		u.broadcastMessage(SocketMessage{"USER_DISCONNECTED", currVid, meta})
	case REQUEST_TIME:
		currVid := p[0]
		u.broadcastMessage(SocketMessage{"SEND_TIME_TO_SERVER", currVid, meta})
	case SEND_TIME_TO_SERVER:
		currVid := p[0]
		currVid.Update(VideoData{currVid.Url, data.Time, true})

		u.broadcastMessage(SocketMessage{"SYNC", currVid, meta})
	case PAUSE_VIDEO:
		if hasVideos {
			currVid := p[0]

			UpdateVideo(u.RoomID, 0, VideoData{currVid.Url, data.Time, false})

			up, err := GetRoomPlaylist(u.RoomID)
			if err != nil {
				log.Logger.Fatal(err)
			}

			if len(up) == 0 {
				u.broadcastMessage(SocketMessage{"PAUSE_VIDEO", VideoData{}, meta})
			} else {
				u.broadcastMessage(SocketMessage{"PAUSE_VIDEO", up[0], meta})
			}
		}
	case SKIP_VIDEO:
		ShiftPlaylistVideo(u.RoomID)

		up, err := GetRoomPlaylist(u.RoomID)
		if err != nil {
			log.Logger.Fatal(err)
		}

		if len(up) == 0 {
			u.broadcastMessage(SocketMessage{"SKIP_VIDEO", VideoData{}, meta})
		} else {
			u.broadcastMessage(SocketMessage{"SKIP_VIDEO", up[0], meta})
		}
	case SYNC:
		currVid := p[0]
		u.broadcastMessage(SocketMessage{"SYNC", currVid, meta})
	default:
		logger.Printf("No valid action sent from Client, ACTION: %v \n", eventType)
	}

	logger.WithFields(logrus.Fields{
		"action":        eventType,
		"data":          data,
		"curr_playlist": p,
	}).Info("[EVENT_HANDLER] Event log")
}

// Unpacks the marsheled json data by the socket message
func unmarshalSocketMessage(msg []byte) (EventType, VideoData) {
	var objmap map[string]json.RawMessage
	err := json.Unmarshal(msg, &objmap)
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	var action EventType
	err = json.Unmarshal(objmap["action"], &action)
	if err != nil {
		log.Logger.Fatal(err)
	}

	var data VideoData
	if objmap["data"] != nil {
		err = json.Unmarshal(objmap["data"], &data)
		if err != nil {
			log.Logger.Fatal(err)
		}
	}

	return action, data
}
