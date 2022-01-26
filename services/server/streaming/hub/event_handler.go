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

	eventType, data := unmarshalSocketMessage(rawMsg)
	itemsInPlaylist := len(Instance.RoomsPlaylist[u.RoomID]) > 0

	meta := EventMetaData{u.Name, GetRoomPop(u.RoomID), u.RoomID}

	switch eventType {
	case REQUEST:
		// TODO: Proper validate if it's a valid youtube video
		Instance.RoomsPlaylist[u.RoomID] = Instance.RoomsPlaylist[u.RoomID].Enqueue(VideoData{
			Time:    0,
			Playing: true,
			Url:     data.Url,
		})

		currVid := Instance.RoomsPlaylist[u.RoomID].GetCurrent()
		// TODO: Validate if END_VIDEO should be sent here?
		u.broadcastMessage(SocketMessage{"END_VIDEO", currVid, meta})

	case END_VIDEO:
		Instance.RoomsPlaylist[u.RoomID] = Instance.RoomsPlaylist[u.RoomID].Dequeue()
		currVid := Instance.RoomsPlaylist[u.RoomID].GetCurrent()

		u.broadcastMessage(SocketMessage{"END_VIDEO", currVid, meta})
	case PLAY_VIDEO:
		if itemsInPlaylist {
			currVid := Instance.RoomsPlaylist[u.RoomID].GetCurrent()
			currVid.Update(VideoData{currVid.Url, data.Time, true})

			u.broadcastMessage(SocketMessage{"PLAY_VIDEO", currVid, meta})
		}
	case USER_JOINED:
		currVid := Instance.RoomsPlaylist[u.RoomID].GetCurrent()
		u.broadcastMessage(SocketMessage{"USER_JOINED", currVid, meta})
	case USER_DISCONNECTED:
		currVid := Instance.RoomsPlaylist[u.RoomID].GetCurrent()
		u.broadcastMessage(SocketMessage{"USER_DISCONNECTED", currVid, meta})
	case REQUEST_TIME:
		currVid := Instance.RoomsPlaylist[u.RoomID].GetCurrent()
		u.broadcastMessage(SocketMessage{"SEND_TIME_TO_SERVER", currVid, meta})
	case SEND_TIME_TO_SERVER:
		currVid := Instance.RoomsPlaylist[u.RoomID].GetCurrent()
		currVid.Update(VideoData{currVid.Url, data.Time, true})

		u.broadcastMessage(SocketMessage{"SYNC", currVid, meta})
	case PAUSE_VIDEO:
		if itemsInPlaylist {
			currVid := Instance.RoomsPlaylist[u.RoomID].GetCurrent()
			currVid.Update(VideoData{currVid.Url, data.Time, false})

			u.broadcastMessage(SocketMessage{"PAUSE_VIDEO", currVid, meta})
		}
	case SYNC:
		currVid := Instance.RoomsPlaylist[u.RoomID].GetCurrent()
		u.broadcastMessage(SocketMessage{"SYNC", currVid, meta})
	default:
		logger.Printf("No valid action sent from Client, ACTION: %v \n", eventType)
	}

	logger.WithFields(logrus.Fields{
		"action":        eventType,
		"data":          data,
		"curr_playlist": Instance.RoomsPlaylist[u.RoomID],
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
