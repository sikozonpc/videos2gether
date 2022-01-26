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
	REQUEST_TIME        EventType = "REQUEST_TIME"
	SEND_TIME_TO_SERVER EventType = "SEND_TIME_TO_SERVER"
)

type EventMetaData struct {
	ActionFrom       string
	UsersConnected int
	RoomID         string
}

func HandleActionEvent(rawMsg []byte, u *User) {
	logger := log.Logger.WithFields(logrus.Fields{
		"room": u.RoomID,
		"user": u.Name,
	})

	eventType, data := unmarshalSocketMessage(rawMsg)
	itemsInPlaylist := len(Instance.RoomsPlaylist[u.RoomID]) > 0

	meta := EventMetaData{u.Name, Instance.getRoomPop(u.RoomID), u.RoomID}

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
		res := SocketMessage{"END_VIDEO", currVid, meta}
		jsData, _ := json.Marshal(res)

		m := Message{jsData, u.RoomID}
		Instance.Broadcast <- m

	case END_VIDEO:
		if len(Instance.RoomsPlaylist[u.RoomID]) > 0 {
			currVid := Instance.RoomsPlaylist[u.RoomID].GetCurrent()
			Instance.RoomsPlaylist[u.RoomID] = Instance.RoomsPlaylist[u.RoomID].Unqueue()

			res := SocketMessage{"END_VIDEO", currVid, meta}
			jsData, _ := json.Marshal(res)

			m := Message{jsData, u.RoomID}
			Instance.Broadcast <- m
		}
	case PLAY_VIDEO:
		if itemsInPlaylist {
			currVid := Instance.RoomsPlaylist[u.RoomID].GetCurrent()
			currVid.Update(VideoData{currVid.Url, data.Time, true})

			res := SocketMessage{"PLAY_VIDEO", currVid, meta}
			jsData, _ := json.Marshal(res)

			m := Message{jsData, u.RoomID}
			Instance.Broadcast <- m
		}
	case REQUEST_TIME:
		res := SocketMessage{
			Action: "SEND_TIME_TO_SERVER",
		}

		jsData, _ := json.Marshal(res)

		m := Message{jsData, u.RoomID}
		Instance.Broadcast <- m
	case SEND_TIME_TO_SERVER:
		currVid := Instance.RoomsPlaylist[u.RoomID].GetCurrent()
		currVid.Update(VideoData{currVid.Url, data.Time, true})

		res := SocketMessage{"SYNC", currVid, meta}
		jsData, _ := json.Marshal(res)

		m := Message{jsData, u.RoomID}
		Instance.Broadcast <- m

	case PAUSE_VIDEO:
		if itemsInPlaylist {
			currVid := Instance.RoomsPlaylist[u.RoomID].GetCurrent()
			currVid.Update(VideoData{currVid.Url, data.Time, false})

			res := SocketMessage{"PAUSE_VIDEO", currVid, meta}
			jsData, _ := json.Marshal(res)

			m := Message{jsData, u.RoomID}
			Instance.Broadcast <- m
		}
	case SYNC:
		currVid := Instance.RoomsPlaylist[u.RoomID].GetCurrent()

		res := SocketMessage{"SYNC", currVid, meta}
		jsData, _ := json.Marshal(res)

		m := Message{jsData, u.RoomID}
		Instance.Broadcast <- m
	default:
		logger.Printf("No valid action sent from Client, ACTION: %v \n", eventType)
	}

	logger.WithFields(logrus.Fields{
		"action":        eventType,
		"data":          data,
		"curr_playlist": Instance.RoomsPlaylist[u.RoomID],
	}).Info("[ROOM]")
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
