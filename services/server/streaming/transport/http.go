package transport

import (
	"fmt"
	"net/http"
	"streamserver/env"
	"streamserver/responses"
	"streamserver/streaming"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type HTTP struct {
	svc streaming.Service
}

// NewHTTP Creates a new HTTP connection
func NewHTTP(svc streaming.Service, r *mux.Router) {
	h := HTTP{svc}

	r.HandleFunc("/room", h.handleCreateRoom).Methods("GET")
	r.HandleFunc("/rooms", h.handleDeleteAllRooms).Methods("DELETE")
	r.HandleFunc("/room/{roomID}", h.handleDeleteRoom).Methods("DELETE")
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

func (h *HTTP) handleDeleteRoom(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	roomId := params["roomID"]

	_, err := h.svc.GetRoomPlaylist(roomId)
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, fmt.Errorf("room does not exist"))
		return
	}

	err = h.svc.DeleteRoom(roomId)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, fmt.Errorf("room cannot be deleted"))
		return
	}

	responses.JSON(w, http.StatusOK, "ok")
}

func (h *HTTP) handleDeleteAllRooms(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("x-api")
	err := checkAPIKey(apiKey)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	h.svc.CleanAllRooms()
	responses.JSON(w, http.StatusOK, "ok")
}

func (h *HTTP) handleGetRoomPlaylist(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	roomID := params["roomID"]

	if len(roomID) == 0 {
		responses.ERROR(w, http.StatusBadRequest, fmt.Errorf("missing roomID param"))
		return
	}

	playlist, err := h.svc.GetRoomPlaylist(roomID)
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, fmt.Errorf("room does not exist"))
		return
	}

	responses.JSON(w, http.StatusOK, playlist)
}

func checkAPIKey(key string) error {
	if len(key) == 0 {
		return fmt.Errorf("missing auth api key")
	}
	if env.Vars.APIKey != key {
		return fmt.Errorf("invalid auth api key")
	}

	return nil
}
