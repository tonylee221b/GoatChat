package in

import (
	"log/slog"
	"net/http"

	"GoatChat/GoatChat/internal/chat/application"
	"GoatChat/GoatChat/internal/chat/domain"
	jsonutil "GoatChat/GoatChat/internal/shared/json_util"

	"github.com/google/uuid"
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

type ChatroomDeleteRequest struct {
	RoomID string `json:"room_id"`
}

type ChatroomUpdateRequest struct {
	RoomID      string `json:"room_id"`
	RoomType    string `json:"room_type"`
	RoomName    string `json:"room_name"`
	Description string `json:"description,omitempty"`
}

func NewChatHandler(crSvc application.ChatService) *ChatHandler {
	return &ChatHandler{crSvc}
}

const (
	ErrDecodeReq  = "failed to decode request"
	ErrInvalidRId = "invalid room id"
	ErrInvalidRt  = "invalid room type"
	ErrInvalidRn  = "invalid room name"
	ErrInvalidRd  = "invalid room description"
	ErrInvalidOID = "invalid owner id"
	ErrCreateCr   = "failed to create chat room"
	ErrDeleteCr   = "failed to delete chat room"
	ErrUpdateCr   = "failed to update chat room"
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

func (h *ChatHandler) DeleteChatroom(w http.ResponseWriter, r *http.Request) {
	var req ChatroomDeleteRequest
	err := jsonutil.ReadJSON(w, r, &req)
	if err != nil {
		slog.Info(ErrDecodeReq)
		jsonutil.WriteError(w, http.StatusBadRequest, ErrDecodeReq)
		return
	}

	id, err := uuid.Parse(req.RoomID)
	if err != nil {
		slog.Info(ErrInvalidRId)
		jsonutil.WriteError(w, http.StatusBadRequest, ErrInvalidRId)
		return
	}

	err = h.crSvc.DeleteChatroom(r.Context(), id)
	if err != nil {
		slog.Error(ErrDeleteCr)
		jsonutil.WriteError(w, http.StatusBadRequest, ErrDeleteCr)
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, "")
}

func (h *ChatHandler) UpdateChatroom(w http.ResponseWriter, r *http.Request) {
	var req ChatroomUpdateRequest
	err := jsonutil.ReadJSON(w, r, &req)
	if err != nil {
		slog.Info(ErrDecodeReq)
		jsonutil.WriteError(w, http.StatusBadRequest, ErrDecodeReq)
		return
	}

	id, err := uuid.Parse(req.RoomID)
	if err != nil {
		slog.Info(ErrInvalidRId)
		jsonutil.WriteError(w, http.StatusBadRequest, ErrInvalidRId)
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

	err = h.crSvc.UpdateChatroom(r.Context(), id, rt, rn, rd)
	if err != nil {
		slog.Error(ErrUpdateCr)
		jsonutil.WriteError(w, http.StatusBadRequest, ErrUpdateCr)
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, "")
}
