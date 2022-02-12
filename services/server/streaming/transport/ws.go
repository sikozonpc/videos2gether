package transport

import (
	"fmt"
	"log"
	"net/http"
	"streamserver/responses"
	"streamserver/streaming"
	"streamserver/streaming/hub"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/gorilla/websocket"
)

// WS represents the web streaming connection
type WS struct {
	svc streaming.Service
}

// Upgrader specifies parameters for upgrading an HTTP connection to a WebSocket connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// NewWS creates a new websocket connection
func NewWS(svc streaming.Service, r *mux.Router) {
	h := WS{svc}
	r.HandleFunc("/ws/{roomID}", h.handleRoomConn).Methods("GET")
}

func (h *WS) handleRoomConn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)
	roomID := params["roomID"]

	if len(roomID) == 0 {
		responses.ERROR(w, http.StatusBadRequest, fmt.Errorf("missing roomID param"))
		return
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}

	c := &hub.Connection{Send: make(chan []byte, 256), WS: ws}

	rnd, _ := uuid.NewRandom()
	newUser := hub.User{
		Conn: c,
		Room: hub.Room{Id: roomID},
		Name: rnd.String(),
	}
	newUser.Connect()
}
