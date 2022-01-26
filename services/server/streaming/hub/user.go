package hub

import (
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

func (u *User) Register() {
	Instance.Register <- *u

	go u.StartSending()
	go u.StartReading()
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
		Instance.Unregister <- *u
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

func (u *User) syncToRoom() {
	currVideo := Instance.RoomsPlaylist[u.RoomID].GetCurrent()

	msg := SocketMessage{
		Action: "SYNC",
		Data:   currVideo,
	}

	u.Conn.WS.WriteJSON(msg)
}
