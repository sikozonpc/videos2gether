package hub

import (
	"encoding/json"
	"streamserver/log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type User struct {
	Name   string
	RoomID string

	Conn *Connection
}

func (u *User) Connect() {
	Instance.Connect <- *u

	go u.StartSending()
	go u.StartReading()

	u.syncToRoom()
}

func (u *User) disconnect() {
	close(u.Conn.Send)
}

func (u *User) StartReading() {
	c := u.Conn

	logger := log.Logger.WithFields(logrus.Fields{
		"room": u.RoomID,
	})

	defer func() {
		meta := EventMetaData{u.Name, GetRoomPop(u.RoomID) - 1, u.RoomID}
		u.broadcastMessage(SocketMessage{"USER_DISCONNECTED", VideoData{}, meta})

		Instance.Disconnect <- *u
		c.WS.Close()
	}()

	c.WS.SetReadLimit(maxMessageSize)
	timeAllowedToRead := time.Now().Add(pongWait)
	c.WS.SetReadDeadline(timeAllowedToRead)
	c.WS.SetPongHandler(func(string) error {
		c.WS.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msg, err := c.WS.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				logger.Fatalf("error: %v", err)
			}
			break
		}

		HandleActionEvent(msg, u)
	}
}

func (u *User) StartSending() {
	c := u.Conn
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.WS.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.Write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (u *User) broadcastMessage(msg SocketMessage) {
	marshalled, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	m := Message{marshalled, u.RoomID}
	Instance.Broadcast <- m
}

// Sync the user to the user with all of the initial data needed
// Should be called after he joins a channel.
func (u *User) syncToRoom() {
	currVideo, _ := GetCurrentVideo(u.RoomID)
	meta := EventMetaData{u.Name, GetRoomPop(u.RoomID) + 1, u.RoomID}

	msg := SocketMessage{"SYNC", currVideo, meta}
	u.Conn.WS.WriteJSON(msg)

	u.broadcastMessage(SocketMessage{"USER_JOINED", currVideo, meta})
}
