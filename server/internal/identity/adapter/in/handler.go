package in

import (
	"log/slog"
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

type RegisterUserReq struct {
	Username string `json:"username"`
}

type UpdateUserContactReq struct {
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

// Response DTO

type UserResponse struct {
	UserID      string     `json:"user_id"`
	Username    string     `json:"username"`
	Email       string     `json:"email,omitempty"`
	PhoneNumber string     `json:"phone_number,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

func fromDomain(user domain.User) UserResponse {
	return UserResponse{
		UserID:      user.ID.String(),
		Username:    user.Username.Value,
		Email:       user.Contact.Email.Value,
		PhoneNumber: user.Contact.PhoneNumber.Value,
		CreatedAt:   user.Audit.CreatedAt,
		UpdatedAt:   user.Audit.UpdatedAt,
		DeletedAt:   user.Audit.DeletedAt,
	}
}

// End DTO

type IdentityHandler struct {
	userSvc *application.UserService
}

func NewIdentityHandler(us *application.UserService) *IdentityHandler {
	return &IdentityHandler{us}
}

func (h *IdentityHandler) Register(w http.ResponseWriter, r *http.Request) {
	req := RegisterUserReq{}

	if err := jsonutil.ReadJSON(w, r, &req); err != nil {
		slog.Info("bad request", "input", req)
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

	jsonutil.WriteJSON(w, http.StatusOK, fromDomain(user))
}

func (h *IdentityHandler) UpdateContact(w http.ResponseWriter, r *http.Request) {
	usernameParam := chi.URLParam(r, "username")
	req := UpdateUserContactReq{}
	if err := jsonutil.ReadJSON(w, r, &req); err != nil {
		jsonutil.WriteError(w, http.StatusBadRequest, err.Error())
	}
	slog.Debug("request parsed", "req", req)

	username, err := domain.NewUsername(usernameParam)
	if err != nil {
		jsonutil.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	slog.Debug("username created", "username", username)

	contact, err := domain.NewContact(req.PhoneNumber, req.Email)
	if err != nil {
		jsonutil.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	slog.Debug("contact created", "contact", contact)

	updated, err := h.userSvc.UpdateContact(r.Context(), username, contact)
	if err != nil {
		jsonutil.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, fromDomain(updated))
}
