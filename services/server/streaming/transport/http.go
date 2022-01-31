package transport

import (
	"fmt"
	"net/http"
	"streamserver/responses"
	"streamserver/streaming"
	"streamserver/streaming/hub"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type HTTP struct {
	svc streaming.Service
}

// NewHTTP Creates a new HTTP connection
func NewHTTP(svc streaming.Service, r *mux.Router) {
	h := HTTP{svc}

	r.HandleFunc("/health", h.getHealth).Methods("GET")
	r.HandleFunc("/room", h.handleCreateRoom).Methods("GET")
	r.HandleFunc("/room/{roomID}/playlist", h.handleGetRoomPlaylist).Methods("GET")
}

func (h *HTTP) handleCreateRoom(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.NewUUID()
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	rd, err := h.svc.CreateRoom(id.String())
	if err != nil {
		responses.ERROR(w, http.StatusConflict, err)
		return
	}

	if rd.ID == "" {
		responses.JSON(w, http.StatusConflict, rd)
		return
	}

	responses.JSON(w, http.StatusOK, rd)
}

func (h *HTTP) getHealth(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "ok")
}

func (h *HTTP) handleGetRoomPlaylist(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	roomID := params["roomID"]

	if len(roomID) == 0 {
		responses.ERROR(w, http.StatusBadRequest, fmt.Errorf("missing roomID param"))
		return
	}

	roomExists := hub.CheckRoomAvailability(roomID)
	if !roomExists {
		responses.ERROR(w, http.StatusBadRequest, fmt.Errorf("room does not exist"))
		return
	}

	playlist := h.svc.GetRoomPlaylist(roomID)

	responses.JSON(w, http.StatusOK, playlist)
}
