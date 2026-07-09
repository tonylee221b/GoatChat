package in

import (
	"GoatChat/GoatChat/internal/chat/application/service"
	"GoatChat/GoatChat/internal/chat/domain"
	jsonutil "GoatChat/GoatChat/internal/shared/json_util"
	"log/slog"
	"net/http"
)

type ChatHandler struct {
	crSvc service.ChatService
}

type ChatroomCreateRequest struct {
	RoomType    string `json:"room_type"`
	RoomName    string `json:"room_name"`
	OwnerId     string `json:"owner_id"`
	Description string `json:"description,omitempty"`
}

type ChatroomCreateResponse struct {
	RoomId string `json:"room_id"`
	// ...
}

func NewChatHandler(crSvc service.ChatService) (*ChatHandler, error) {
	return &ChatHandler{crSvc}, nil
}

const (
	ErrDecodeReq  = "failed to decode request"
	ErrInvalidRt  = "invalid room type"
	ErrInvalidRn  = "invalid room name"
	ErrInvalidRd  = "invalid room description"
	ErrInvalidOid = "invalid owner id"
	ErrCreateCr   = "failed to create chat room"
)

func (h *ChatHandler) CreateChatroom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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

	oid, err := domain.NewRoomOwnerId(req.OwnerId)
	if err != nil {
		slog.Info(ErrInvalidOid)
		jsonutil.WriteError(w, http.StatusBadRequest, ErrInvalidOid)
		return
	}

	room, err := h.crSvc.CreateChatroom(r.Context(),
		rt,
		rn,
		rd,
		oid)

	if err != nil {
		slog.Error(ErrCreateCr)
		jsonutil.WriteError(w, http.StatusBadRequest, ErrCreateCr)
		return
	}

	res := ChatroomCreateResponse{
		RoomId: room.ID.String(),
	}

	jsonutil.WriteJSON(w, http.StatusCreated, res)
}
