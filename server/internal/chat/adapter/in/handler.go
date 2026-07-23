package in

import (
	"log/slog"
	"net/http"

	"GoatChat/GoatChat/internal/chat/application"
	"GoatChat/GoatChat/internal/chat/domain"
	jsonutil "GoatChat/GoatChat/internal/shared/json_util"
)

type ChatHandler struct {
	crSvc application.ChatService
}

type ChatroomCreateRequest struct {
	RoomType    string `json:"room_type"`
	RoomName    string `json:"room_name"`
	OwnerID     string `json:"owner_id"`
	Description string `json:"description,omitempty"`
}

type ChatroomCreateResponse struct {
	RoomID string `json:"room_id"`
	// ...
}

func NewChatHandler(crSvc application.ChatService) *ChatHandler {
	return &ChatHandler{crSvc}
}

const (
	ErrDecodeReq  = "failed to decode request"
	ErrInvalidRt  = "invalid room type"
	ErrInvalidRn  = "invalid room name"
	ErrInvalidRd  = "invalid room description"
	ErrInvalidOID = "invalid owner id"
	ErrCreateCr   = "failed to create chat room"
)

func (h *ChatHandler) CreateChatroom(w http.ResponseWriter, r *http.Request) {
	var req ChatroomCreateRequest
	err := jsonutil.ReadJSON(w, r, &req)
	if err != nil {
		slog.Info(ErrDecodeReq)
		jsonutil.WriteError(w, http.StatusBadRequest, ErrDecodeReq)
		return
	}

	rt, err := domain.NewRoomType(req.RoomType)
	if err != nil {
		slog.Info(ErrInvalidRt)
		jsonutil.WriteError(w, http.StatusBadRequest, ErrInvalidRt)
		return
	}

	rn, err := domain.NewRoomName(req.RoomName)
	if err != nil {
		slog.Info(ErrInvalidRn)
		jsonutil.WriteError(w, http.StatusBadRequest, ErrInvalidRn)
		return
	}

	rd, err := domain.NewRoomDescription(req.Description)
	if err != nil {
		slog.Info(ErrInvalidRd)
		jsonutil.WriteError(w, http.StatusBadRequest, ErrInvalidRd)
		return
	}

	oid, err := domain.NewRoomOwnerId(req.OwnerID)
	if err != nil {
		slog.Info(ErrInvalidOID)
		jsonutil.WriteError(w, http.StatusBadRequest, ErrInvalidOID)
		return
	}

	room, err := h.crSvc.CreateChatroom(r.Context(), rt, rn, rd, oid)
	if err != nil {
		slog.Error(ErrCreateCr)
		jsonutil.WriteError(w, http.StatusBadRequest, ErrCreateCr)
		return
	}

	res := ChatroomCreateResponse{
		RoomID: room.ID.String(),
	}

	jsonutil.WriteJSON(w, http.StatusCreated, res)
}
