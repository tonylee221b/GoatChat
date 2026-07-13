package in

import (
	"net/http"
	"time"

	"GoatChat/GoatChat/internal/identity/application"
	"GoatChat/GoatChat/internal/identity/domain"
	jsonutil "GoatChat/GoatChat/internal/shared/json_util"

	"github.com/go-chi/chi/v5"
)

const (
	ErrInvalidUsername = "invalid username"
)

// Request DTO

type UserRegisterRequest struct {
	Username string `json:"username"`
}

// Response DTO

type UserResponse struct {
	UserID      string     `json:"user_id"`
	Username    string     `json:"username"`
	Email       string     `json:"email,omitempty"`
	PhoneNumber string     `json:"phone_number,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

// End DTO

type IdentityHandler struct {
	userSvc *application.UserService
}

func NewIdentityHandler(us *application.UserService) *IdentityHandler {
	return &IdentityHandler{us}
}

func (h *IdentityHandler) Register(w http.ResponseWriter, r *http.Request) {
	req := UserRegisterRequest{}

	if err := jsonutil.ReadJSON(w, r, &req); err != nil {
		jsonutil.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	username, err := domain.NewUsername(req.Username)
	if err != nil {
		jsonutil.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.userSvc.Register(r.Context(), username)
	if err != nil {
		jsonutil.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	jsonutil.WriteJSON(w, http.StatusCreated, "user registered successfully")
}

func (h *IdentityHandler) FindByUsername(w http.ResponseWriter, r *http.Request) {
	usernameParam := chi.URLParam(r, "username")
	username, err := domain.NewUsername(usernameParam)
	if err != nil {
		jsonutil.WriteError(w, http.StatusBadRequest, ErrInvalidUsername)
		return
	}

	user, err := h.userSvc.FindByUsername(r.Context(), username)
	if err != nil {
		jsonutil.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, user)
}
