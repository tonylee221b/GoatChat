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

const (
	refreshTokenCookieKey = "refresh_token"
	authPath              = "/api/v1/auth"
)

// Request DTO

type RegisterUserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUserContactReq struct {
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
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

func newUserResponse(user domain.User) UserResponse {
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

type LoginResponse struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// End DTO

type IdentityHandler struct {
	userSvc *application.UserService
	authSvc *application.AuthService
}

func NewIdentityHandler(us *application.UserService, as *application.AuthService) *IdentityHandler {
	return &IdentityHandler{us, as}
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

	err = h.userSvc.Register(r.Context(), username, req.Password)
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

	jsonutil.WriteJSON(w, http.StatusOK, newUserResponse(user))
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

	jsonutil.WriteJSON(w, http.StatusOK, newUserResponse(updated))
}

func (h *IdentityHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq LoginReq
	if err := jsonutil.ReadJSON(w, r, &loginReq); err != nil {
		jsonutil.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	slog.Debug("login request parsed", "username", loginReq.Username)

	username, err := domain.NewUsername(loginReq.Username)
	if err != nil {
		jsonutil.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	tokenPair, err := h.authSvc.Login(r.Context(), username, loginReq.Password)
	if err != nil {
		jsonutil.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	h.setSecureCookie(w, tokenPair.RefreshToken, tokenPair.RefreshTokenExpiresAt)

	loginResp := LoginResponse{
		AccessToken: tokenPair.AccessToken,
		TokenType:   "Bearer",
		ExpiresAt:   tokenPair.AccessTokenExpiresAt,
	}
	jsonutil.WriteJSON(w, http.StatusOK, loginResp)
}

func (h *IdentityHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(refreshTokenCookieKey)
	if err == nil {
		if err := h.authSvc.Logout(r.Context(), cookie.Value); err != nil {
			slog.Warn("logout failed. manual check is required")
			jsonutil.WriteError(w, http.StatusInternalServerError, "logout failed")
			return
		}
	}

	h.clearCookie(w)

	w.WriteHeader(http.StatusNoContent)
}

func (h *IdentityHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(refreshTokenCookieKey)
	if err != nil {
		jsonutil.WriteError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	tokenPair, err := h.authSvc.Refresh(r.Context(), cookie.Value)
	if err != nil {
		jsonutil.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.setSecureCookie(w, tokenPair.RefreshToken, tokenPair.RefreshTokenExpiresAt)

	loginResp := LoginResponse{
		AccessToken: tokenPair.AccessToken,
		TokenType:   "Bearer",
		ExpiresAt:   tokenPair.AccessTokenExpiresAt,
	}
	jsonutil.WriteJSON(w, http.StatusOK, loginResp)
}

func (h *IdentityHandler) setSecureCookie(w http.ResponseWriter, refreshToken string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshTokenCookieKey,
		Value:    refreshToken,
		Path:     authPath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  expiresAt,
		MaxAge:   int(time.Until(expiresAt).Seconds()),
	})
}

func (h *IdentityHandler) clearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshTokenCookieKey,
		Value:    "",
		Path:     authPath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})
}
