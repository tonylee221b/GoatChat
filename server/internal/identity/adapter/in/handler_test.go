package in_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"GoatChat/GoatChat/internal/identity/adapter/in"
	"GoatChat/GoatChat/internal/identity/application"
	"GoatChat/GoatChat/internal/identity/application/port/mocks"
	"GoatChat/GoatChat/internal/identity/domain"
	"GoatChat/GoatChat/internal/shared/db_tx/testutils"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRegisterUser(t *testing.T) {
	for _, tt := range registerUserTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			require.NoError(t, json.NewEncoder(&buf).Encode(tt.body))

			mUserRepo := mocks.NewMockUserRepository(t)
			mSessionRepo := mocks.NewMockSessionRepository(t)
			mPWHasher := mocks.NewMockPasswordHasher(t)
			mAccessTokenManager := mocks.NewMockAccessTokenManager(t)
			if tt.setupMock != nil {
				tt.setupMock(mUserRepo, mPWHasher)
			}

			userSvc := application.NewUserService(testutils.StubTx{}, mUserRepo, mPWHasher)
			authSvc := application.NewAuthService(testutils.StubTx{}, mUserRepo, mSessionRepo, mPWHasher, mAccessTokenManager, testAccessTokenTTL, testRefreshTokenTTL)
			h := in.NewIdentityHandler(userSvc, authSvc)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/mock", &buf)
			req.Header.Set("Content-Type", "application/json")
			h.Register(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			require.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestFindUserByUsername(t *testing.T) {
	for _, tt := range findUserTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			mUserRepo := mocks.NewMockUserRepository(t)
			mSessionRepo := mocks.NewMockSessionRepository(t)
			mPWHasher := mocks.NewMockPasswordHasher(t)
			mAccessTokenManager := mocks.NewMockAccessTokenManager(t)
			if tt.setupMock != nil {
				tt.setupMock(mUserRepo)
			}

			userSvc := application.NewUserService(testutils.StubTx{}, mUserRepo, mPWHasher)
			authSvc := application.NewAuthService(testutils.StubTx{}, mUserRepo, mSessionRepo, mPWHasher, mAccessTokenManager, testAccessTokenTTL, testRefreshTokenTTL)
			h := in.NewIdentityHandler(userSvc, authSvc)

			r := chi.NewRouter()
			r.Get("/users/{username}", h.FindByUsername)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/users/"+tt.param, nil)
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			require.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestUpdateUserContact(t *testing.T) {
	for _, tt := range updateContactTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			mUserRepo := mocks.NewMockUserRepository(t)
			mSessionRepo := mocks.NewMockSessionRepository(t)
			mPWHasher := mocks.NewMockPasswordHasher(t)
			mAccessTokenManager := mocks.NewMockAccessTokenManager(t)
			if tt.setupMock != nil {
				tt.setupMock(mUserRepo)
			}

			userSvc := application.NewUserService(testutils.StubTx{}, mUserRepo, mPWHasher)
			authSvc := application.NewAuthService(testutils.StubTx{}, mUserRepo, mSessionRepo, mPWHasher, mAccessTokenManager, testAccessTokenTTL, testRefreshTokenTTL)
			h := in.NewIdentityHandler(userSvc, authSvc)

			r := chi.NewRouter()
			r.Put("/users/{username}", h.UpdateContact)

			body, err := json.Marshal(tt.body)
			require.NoError(t, err)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPut, "/users/"+tt.param, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			require.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestLogin(t *testing.T) {
	for _, tt := range loginTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			mUserRepo := mocks.NewMockUserRepository(t)
			mSessionRepo := mocks.NewMockSessionRepository(t)
			mPWHasher := mocks.NewMockPasswordHasher(t)
			mAccessTokenManager := mocks.NewMockAccessTokenManager(t)
			if tt.setupMock != nil {
				tt.setupMock(mUserRepo, mSessionRepo, mPWHasher, mAccessTokenManager)
			}

			userSvc := application.NewUserService(testutils.StubTx{}, mUserRepo, mPWHasher)
			authSvc := application.NewAuthService(testutils.StubTx{}, mUserRepo, mSessionRepo, mPWHasher, mAccessTokenManager, testAccessTokenTTL, testRefreshTokenTTL)
			h := in.NewIdentityHandler(userSvc, authSvc)

			r := chi.NewRouter()
			r.Post("/auth/login", h.Login)

			body, err := json.Marshal(tt.body)
			require.NoError(t, err)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			require.Equal(t, tt.wantStatus, res.StatusCode)
			if tt.wantStatus != http.StatusOK {
				return
			}

			var loginResp in.LoginResponse
			require.NoError(t, json.NewDecoder(res.Body).Decode(&loginResp))
			require.Equal(t, "access-token", loginResp.AccessToken)

			var refreshCookie *http.Cookie
			for _, cookie := range res.Cookies() {
				if cookie.Name == "refresh_token" {
					refreshCookie = cookie
					break
				}
			}
			require.NotNil(t, refreshCookie)
			require.NotEmpty(t, refreshCookie.Value)
			require.True(t, refreshCookie.HttpOnly)
			require.True(t, refreshCookie.Secure)
			require.Equal(t, "/api/v1/auth", refreshCookie.Path)
			require.Equal(t, http.SameSiteStrictMode, refreshCookie.SameSite)
		})
	}
}

func TestRefresh(t *testing.T) {
	for _, tt := range refreshTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			mUserRepo := mocks.NewMockUserRepository(t)
			mSessionRepo := mocks.NewMockSessionRepository(t)
			mPWHasher := mocks.NewMockPasswordHasher(t)
			mAccessTokenManager := mocks.NewMockAccessTokenManager(t)
			if tt.setupMock != nil {
				tt.setupMock(mSessionRepo, mAccessTokenManager)
			}

			userSvc := application.NewUserService(testutils.StubTx{}, mUserRepo, mPWHasher)
			authSvc := application.NewAuthService(testutils.StubTx{}, mUserRepo, mSessionRepo, mPWHasher, mAccessTokenManager, testAccessTokenTTL, testRefreshTokenTTL)
			h := in.NewIdentityHandler(userSvc, authSvc)

			r := chi.NewRouter()
			r.Post("/auth/refresh", h.Refresh)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
			if tt.withCookie {
				req.AddCookie(&http.Cookie{Name: "refresh_token", Value: tt.refreshToken})
			}

			r.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			require.Equal(t, tt.wantStatus, res.StatusCode)
			if tt.wantStatus != http.StatusOK {
				require.Empty(t, res.Cookies())
				return
			}

			var refreshResp in.LoginResponse
			require.NoError(t, json.NewDecoder(res.Body).Decode(&refreshResp))
			require.Equal(t, tt.wantAccessToken, refreshResp.AccessToken)

			var refreshCookie *http.Cookie
			for _, cookie := range res.Cookies() {
				if cookie.Name == "refresh_token" {
					refreshCookie = cookie
					break
				}
			}
			require.NotNil(t, refreshCookie)
			require.NotEmpty(t, refreshCookie.Value)
			require.NotEqual(t, tt.refreshToken, refreshCookie.Value)
			require.True(t, refreshCookie.HttpOnly)
			require.True(t, refreshCookie.Secure)
			require.Equal(t, "/api/v1/auth", refreshCookie.Path)
			require.Equal(t, http.SameSiteStrictMode, refreshCookie.SameSite)
		})
	}
}

func TestLogout(t *testing.T) {
	for _, tt := range logoutTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			mUserRepo := mocks.NewMockUserRepository(t)
			mSessionRepo := mocks.NewMockSessionRepository(t)
			mPWHasher := mocks.NewMockPasswordHasher(t)
			mAccessTokenManager := mocks.NewMockAccessTokenManager(t)
			if tt.setupMock != nil {
				tt.setupMock(mSessionRepo)
			}

			userSvc := application.NewUserService(testutils.StubTx{}, mUserRepo, mPWHasher)
			authSvc := application.NewAuthService(testutils.StubTx{}, mUserRepo, mSessionRepo, mPWHasher, mAccessTokenManager, testAccessTokenTTL, testRefreshTokenTTL)
			h := in.NewIdentityHandler(userSvc, authSvc)

			r := chi.NewRouter()
			r.Post("/auth/logout", h.Logout)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
			if tt.withCookie {
				req.AddCookie(&http.Cookie{Name: "refresh_token", Value: tt.refreshToken})
			}

			r.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			require.Equal(t, tt.wantStatus, res.StatusCode)
			if tt.wantStatus != http.StatusNoContent {
				require.Empty(t, res.Cookies())
				return
			}

			var refreshCookie *http.Cookie
			for _, cookie := range res.Cookies() {
				if cookie.Name == "refresh_token" {
					refreshCookie = cookie
					break
				}
			}
			require.NotNil(t, refreshCookie)
			require.Empty(t, refreshCookie.Value)
			require.Equal(t, -1, refreshCookie.MaxAge)
			require.True(t, refreshCookie.HttpOnly)
			require.True(t, refreshCookie.Secure)
			require.Equal(t, "/api/v1/auth", refreshCookie.Path)
			require.Equal(t, http.SameSiteStrictMode, refreshCookie.SameSite)
		})
	}
}

const (
	testPlainPassword   = "plain-password"
	testAccessTokenTTL  = 15 * time.Minute
	testRefreshTokenTTL = 7 * 24 * time.Hour
)

type registerUserTestCase struct {
	name      string
	body      in.RegisterUserReq
	setupMock func(
		*mocks.MockUserRepository,
		*mocks.MockPasswordHasher,
	)
	wantStatus int
}

func registerUserTestCases() []registerUserTestCase {
	username, _ := domain.NewUsername("test-user")
	hashedPassword := "hashed-password"
	passwordHash, _ := domain.NewPasswordHash(hashedPassword)

	findUserErr := errors.New(domain.ErrDB)
	hashPasswordErr := errors.New("failed to hash password")
	saveUserErr := errors.New(domain.ErrDB)

	return []registerUserTestCase{
		{
			name: "success",
			body: in.RegisterUserReq{
				Username: username.Value,
				Password: "plain-password",
			},
			setupMock: func(
				userRepo *mocks.MockUserRepository,
				pwHasher *mocks.MockPasswordHasher,
			) {
				userRepo.EXPECT().
					ExistsByUsername(mock.Anything, username).
					Return(false, nil).
					Once()

				pwHasher.EXPECT().
					Hash("plain-password").
					Return(hashedPassword, nil).
					Once()

				userRepo.EXPECT().
					Save(mock.Anything, mock.MatchedBy(func(user domain.User) bool {
						return user.Username == username &&
							user.PasswordHash == passwordHash
					})).
					Return(nil).
					Once()
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "invalid username",
			body: in.RegisterUserReq{
				Username: "",
				Password: "plain-password",
			},
			setupMock:  nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "user already exists",
			body: in.RegisterUserReq{
				Username: username.Value,
				Password: "plain-password",
			},
			setupMock: func(
				userRepo *mocks.MockUserRepository,
				_ *mocks.MockPasswordHasher,
			) {
				userRepo.EXPECT().
					ExistsByUsername(mock.Anything, username).
					Return(true, nil).
					Once()
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "user lookup failure",
			body: in.RegisterUserReq{
				Username: username.Value,
				Password: "plain-password",
			},
			setupMock: func(
				userRepo *mocks.MockUserRepository,
				_ *mocks.MockPasswordHasher,
			) {
				userRepo.EXPECT().
					ExistsByUsername(mock.Anything, username).
					Return(false, findUserErr).
					Once()
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "password hash failure",
			body: in.RegisterUserReq{
				Username: username.Value,
				Password: "plain-password",
			},
			setupMock: func(
				userRepo *mocks.MockUserRepository,
				pwHasher *mocks.MockPasswordHasher,
			) {
				userRepo.EXPECT().
					ExistsByUsername(mock.Anything, username).
					Return(false, nil).
					Once()

				pwHasher.EXPECT().
					Hash("plain-password").
					Return("", hashPasswordErr).
					Once()
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid password hash",
			body: in.RegisterUserReq{
				Username: username.Value,
				Password: "plain-password",
			},
			setupMock: func(
				userRepo *mocks.MockUserRepository,
				pwHasher *mocks.MockPasswordHasher,
			) {
				userRepo.EXPECT().
					ExistsByUsername(mock.Anything, username).
					Return(false, nil).
					Once()

				pwHasher.EXPECT().
					Hash("plain-password").
					Return("", nil).
					Once()
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "user save failure",
			body: in.RegisterUserReq{
				Username: username.Value,
				Password: "plain-password",
			},
			setupMock: func(
				userRepo *mocks.MockUserRepository,
				pwHasher *mocks.MockPasswordHasher,
			) {
				userRepo.EXPECT().
					ExistsByUsername(mock.Anything, username).
					Return(false, nil).
					Once()

				pwHasher.EXPECT().
					Hash("plain-password").
					Return(hashedPassword, nil).
					Once()

				userRepo.EXPECT().
					Save(mock.Anything, mock.Anything).
					Return(saveUserErr).
					Once()
			},
			wantStatus: http.StatusBadRequest,
		},
	}
}

type findUserTestCase struct {
	name       string
	param      string
	setupMock  func(*mocks.MockUserRepository)
	wantStatus int
}

func findUserTestCases() []findUserTestCase {
	username := "test-user"
	domainUsername, _ := domain.NewUsername(username)
	passwordHash, _ := domain.NewPasswordHash("hashed-password")
	user := domain.NewUser(domainUsername, passwordHash)

	userNotFound := errors.New(domain.ErrUserNotFound)
	repositoryErr := errors.New(domain.ErrDB)

	return []findUserTestCase{
		{
			name:  "success",
			param: username,
			setupMock: func(userRepo *mocks.MockUserRepository) {
				userRepo.EXPECT().
					FindByUsername(mock.Anything, domainUsername).
					Return(user, nil).
					Once()
			},
			wantStatus: http.StatusOK,
		},
		{
			name:  "user not found",
			param: "wrong-user",
			setupMock: func(userRepo *mocks.MockUserRepository) {
				userRepo.EXPECT().
					FindByUsername(mock.Anything, mock.Anything).
					Return(nil, userNotFound).
					Once()
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:  "repository failure",
			param: username,
			setupMock: func(userRepo *mocks.MockUserRepository) {
				userRepo.EXPECT().
					FindByUsername(mock.Anything, domainUsername).
					Return(nil, repositoryErr).
					Once()
			},
			wantStatus: http.StatusNotFound,
		},
	}
}

type updateContactTestCase struct {
	name       string
	param      string
	body       in.UpdateUserContactReq
	setupMock  func(*mocks.MockUserRepository)
	wantStatus int
}

func updateContactTestCases() []updateContactTestCase {
	username := "test-user"
	domainUsername, _ := domain.NewUsername(username)
	passwordHash, _ := domain.NewPasswordHash("hashed-password")

	successUser := domain.NewUser(domainUsername, passwordHash)
	updateFailureUser := domain.NewUser(domainUsername, passwordHash)

	userNotFound := errors.New(domain.ErrUserNotFound)
	updateUserErr := errors.New(domain.ErrDB)

	validContact := in.UpdateUserContactReq{
		PhoneNumber: "010-1234-5678",
		Email:       "test@example.com",
	}

	return []updateContactTestCase{
		{
			name:  "success",
			param: username,
			body:  validContact,
			setupMock: func(userRepo *mocks.MockUserRepository) {
				userRepo.EXPECT().
					FindByUsername(mock.Anything, domainUsername).
					Return(successUser, nil).
					Once()

				userRepo.EXPECT().
					Update(mock.Anything, mock.MatchedBy(func(user domain.User) bool {
						return user.Contact.PhoneNumber.Value == validContact.PhoneNumber &&
							user.Contact.Email.Value == validContact.Email
					})).
					Return(nil).
					Once()
			},
			wantStatus: http.StatusOK,
		},
		{
			name:  "invalid phone number",
			param: username,
			body: in.UpdateUserContactReq{
				PhoneNumber: "invalid-phone-number",
				Email:       "test@example.com",
			},
			setupMock:  nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:  "invalid email",
			param: username,
			body: in.UpdateUserContactReq{
				PhoneNumber: "010-1234-5678",
				Email:       "",
			},
			setupMock:  nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:  "user not found",
			param: username,
			body:  validContact,
			setupMock: func(userRepo *mocks.MockUserRepository) {
				userRepo.EXPECT().
					FindByUsername(mock.Anything, domainUsername).
					Return(nil, userNotFound).
					Once()
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:  "repository update failure",
			param: username,
			body:  validContact,
			setupMock: func(userRepo *mocks.MockUserRepository) {
				userRepo.EXPECT().
					FindByUsername(mock.Anything, domainUsername).
					Return(updateFailureUser, nil).
					Once()

				userRepo.EXPECT().
					Update(mock.Anything, mock.Anything).
					Return(updateUserErr).
					Once()
			},
			wantStatus: http.StatusBadRequest,
		},
	}
}

type loginTestCase struct {
	name      string
	body      in.LoginReq
	setupMock func(
		*mocks.MockUserRepository,
		*mocks.MockSessionRepository,
		*mocks.MockPasswordHasher,
		*mocks.MockAccessTokenManager,
	)
	wantStatus int
}

func loginTestCases() []loginTestCase {
	username, _ := domain.NewUsername("test-user")
	passwordHash, _ := domain.NewPasswordHash("hashed-password")
	user := domain.NewUser(username, passwordHash)

	const (
		plainPassword = "plain-password"
		accessToken   = "access-token"
	)

	userNotFoundErr := errors.New(domain.ErrUserNotFound)
	repositoryErr := errors.New(domain.ErrDB)
	passwordMismatchErr := errors.New(domain.ErrInvalidCredentials)
	sessionSaveErr := errors.New(domain.ErrDB)

	validBody := in.LoginReq{
		Username: username.Value,
		Password: plainPassword,
	}

	return []loginTestCase{
		{
			name: "success",
			body: validBody,
			setupMock: func(
				userRepo *mocks.MockUserRepository,
				sessionRepo *mocks.MockSessionRepository,
				pwHasher *mocks.MockPasswordHasher,
				accessTokenManager *mocks.MockAccessTokenManager,
			) {
				userRepo.EXPECT().
					FindByUsername(mock.Anything, username).
					Return(user, nil).
					Once()

				pwHasher.EXPECT().
					Compare(passwordHash.Value, plainPassword).
					Return(nil).
					Once()

				accessTokenManager.EXPECT().
					Issue(user.ID.String(), mock.Anything).
					Return(accessToken, nil).
					Once()

				sessionRepo.EXPECT().
					Save(mock.Anything, mock.Anything).
					Return(nil).
					Once()
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "invalid username",
			body: in.LoginReq{
				Username: "",
				Password: plainPassword,
			},
			setupMock:  nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "user not found",
			body: validBody,
			setupMock: func(
				userRepo *mocks.MockUserRepository,
				_ *mocks.MockSessionRepository,
				_ *mocks.MockPasswordHasher,
				_ *mocks.MockAccessTokenManager,
			) {
				userRepo.EXPECT().
					FindByUsername(mock.Anything, username).
					Return(nil, userNotFoundErr).
					Once()
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "repository failure",
			body: validBody,
			setupMock: func(
				userRepo *mocks.MockUserRepository,
				_ *mocks.MockSessionRepository,
				_ *mocks.MockPasswordHasher,
				_ *mocks.MockAccessTokenManager,
			) {
				userRepo.EXPECT().
					FindByUsername(mock.Anything, username).
					Return(nil, repositoryErr).
					Once()
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "password mismatch",
			body: validBody,
			setupMock: func(
				userRepo *mocks.MockUserRepository,
				_ *mocks.MockSessionRepository,
				pwHasher *mocks.MockPasswordHasher,
				_ *mocks.MockAccessTokenManager,
			) {
				userRepo.EXPECT().
					FindByUsername(mock.Anything, username).
					Return(user, nil).
					Once()

				pwHasher.EXPECT().
					Compare(passwordHash.Value, plainPassword).
					Return(passwordMismatchErr).
					Once()
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "session save failure",
			body: validBody,
			setupMock: func(
				userRepo *mocks.MockUserRepository,
				sessionRepo *mocks.MockSessionRepository,
				pwHasher *mocks.MockPasswordHasher,
				accessTokenManager *mocks.MockAccessTokenManager,
			) {
				userRepo.EXPECT().
					FindByUsername(mock.Anything, username).
					Return(user, nil).
					Once()

				pwHasher.EXPECT().
					Compare(passwordHash.Value, plainPassword).
					Return(nil).
					Once()

				accessTokenManager.EXPECT().
					Issue(user.ID.String(), mock.Anything).
					Return(accessToken, nil).
					Once()

				sessionRepo.EXPECT().
					Save(mock.Anything, mock.Anything).
					Return(sessionSaveErr).
					Once()
			},
			wantStatus: http.StatusUnauthorized,
		},
	}
}

type refreshTestCase struct {
	name         string
	refreshToken string
	withCookie   bool
	setupMock    func(
		*mocks.MockSessionRepository,
		*mocks.MockAccessTokenManager,
	)
	wantAccessToken string
	wantStatus      int
}

func refreshTestCases() []refreshTestCase {
	const (
		refreshToken = "refresh-token"
		accessToken  = "access-token"
	)

	currentHash, _ := domain.NewRefreshTokenHash(refreshToken)
	username, _ := domain.NewUsername("test-user")
	passwordHash, _ := domain.NewPasswordHash("hashed-password")
	user := domain.NewUser(username, passwordHash)
	now := time.Now().UTC()

	newSession := func(expiresAt time.Time) *domain.Session {
		session, _ := domain.NewSession(
			user.ID,
			currentHash,
			now.Add(-2*time.Hour),
			expiresAt,
		)
		return session
	}

	successSession := newSession(now.Add(testRefreshTokenTTL))
	revokedSession := newSession(now.Add(testRefreshTokenTTL))
	expiredSession := newSession(now.Add(-time.Minute))
	issueFailureSession := newSession(now.Add(testRefreshTokenTTL))
	updateFailureSession := newSession(now.Add(testRefreshTokenTTL))

	revokedAt := now.Add(-time.Minute)
	revokedSession.RevokedAt = &revokedAt

	return []refreshTestCase{
		{
			name:         "success",
			refreshToken: refreshToken,
			withCookie:   true,
			setupMock: func(
				sessionRepo *mocks.MockSessionRepository,
				accessTokenManager *mocks.MockAccessTokenManager,
			) {
				sessionRepo.EXPECT().
					FindByRefreshTokenHash(mock.Anything, currentHash).
					Return(successSession, nil).
					Once()

				accessTokenManager.EXPECT().
					Issue(successSession.UserID.String(), successSession.ID.String()).
					Return(accessToken, nil).
					Once()

				sessionRepo.EXPECT().
					UpdateRefreshTokenHash(
						mock.Anything,
						successSession.ID,
						mock.MatchedBy(func(newHash domain.RefreshTokenHash) bool {
							return !newHash.IsZero() && !newHash.Equals(currentHash)
						}),
					).
					Return(nil).
					Once()
			},
			wantAccessToken: accessToken,
			wantStatus:      http.StatusOK,
		},
		{
			name:         "missing refresh token cookie",
			refreshToken: "",
			withCookie:   false,
			setupMock:    nil,
			wantStatus:   http.StatusUnauthorized,
		},
		{
			name:         "empty refresh token",
			refreshToken: "",
			withCookie:   true,
			setupMock:    nil,
			wantStatus:   http.StatusBadRequest,
		},
		{
			name:         "session not found",
			refreshToken: refreshToken,
			withCookie:   true,
			setupMock: func(
				sessionRepo *mocks.MockSessionRepository,
				_ *mocks.MockAccessTokenManager,
			) {
				sessionRepo.EXPECT().
					FindByRefreshTokenHash(mock.Anything, currentHash).
					Return(nil, errors.New(domain.ErrSessionNotFound)).
					Once()
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:         "revoked session",
			refreshToken: refreshToken,
			withCookie:   true,
			setupMock: func(
				sessionRepo *mocks.MockSessionRepository,
				_ *mocks.MockAccessTokenManager,
			) {
				sessionRepo.EXPECT().
					FindByRefreshTokenHash(mock.Anything, currentHash).
					Return(revokedSession, nil).
					Once()
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:         "expired session",
			refreshToken: refreshToken,
			withCookie:   true,
			setupMock: func(
				sessionRepo *mocks.MockSessionRepository,
				_ *mocks.MockAccessTokenManager,
			) {
				sessionRepo.EXPECT().
					FindByRefreshTokenHash(mock.Anything, currentHash).
					Return(expiredSession, nil).
					Once()
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:         "access token issue failure",
			refreshToken: refreshToken,
			withCookie:   true,
			setupMock: func(
				sessionRepo *mocks.MockSessionRepository,
				accessTokenManager *mocks.MockAccessTokenManager,
			) {
				sessionRepo.EXPECT().
					FindByRefreshTokenHash(mock.Anything, currentHash).
					Return(issueFailureSession, nil).
					Once()

				accessTokenManager.EXPECT().
					Issue(issueFailureSession.UserID.String(), issueFailureSession.ID.String()).
					Return("", errors.New(domain.ErrTokenUserIDRequired)).
					Once()
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:         "refresh token update failure",
			refreshToken: refreshToken,
			withCookie:   true,
			setupMock: func(
				sessionRepo *mocks.MockSessionRepository,
				accessTokenManager *mocks.MockAccessTokenManager,
			) {
				sessionRepo.EXPECT().
					FindByRefreshTokenHash(mock.Anything, currentHash).
					Return(updateFailureSession, nil).
					Once()

				accessTokenManager.EXPECT().
					Issue(updateFailureSession.UserID.String(), updateFailureSession.ID.String()).
					Return(accessToken, nil).
					Once()

				sessionRepo.EXPECT().
					UpdateRefreshTokenHash(mock.Anything, updateFailureSession.ID, mock.Anything).
					Return(errors.New(domain.ErrDB)).
					Once()
			},
			wantStatus: http.StatusBadRequest,
		},
	}
}

type logoutTestCase struct {
	name         string
	refreshToken string
	withCookie   bool
	setupMock    func(*mocks.MockSessionRepository)
	wantStatus   int
}

func logoutTestCases() []logoutTestCase {
	const refreshToken = "refresh-token"

	refreshTokenHash, _ := domain.NewRefreshTokenHash(refreshToken)
	username, _ := domain.NewUsername("test-user")
	passwordHash, _ := domain.NewPasswordHash("hashed-password")
	user := domain.NewUser(username, passwordHash)
	now := time.Now().UTC()

	newSession := func() *domain.Session {
		session, _ := domain.NewSession(
			user.ID,
			refreshTokenHash,
			now,
			now.Add(testRefreshTokenTTL),
		)
		return session
	}

	successSession := newSession()
	revokeFailureSession := newSession()

	return []logoutTestCase{
		{
			name:         "success",
			refreshToken: refreshToken,
			withCookie:   true,
			setupMock: func(sessionRepo *mocks.MockSessionRepository) {
				sessionRepo.EXPECT().
					FindByRefreshTokenHashForUpdate(mock.Anything, refreshTokenHash).
					Return(successSession, nil).
					Once()

				sessionRepo.EXPECT().
					Revoke(mock.Anything, successSession.ID, mock.Anything).
					Return(nil).
					Once()
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:         "missing refresh token cookie",
			refreshToken: "",
			withCookie:   false,
			setupMock:    nil,
			wantStatus:   http.StatusNoContent,
		},
		{
			name:         "empty refresh token",
			refreshToken: "",
			withCookie:   true,
			setupMock:    nil,
			wantStatus:   http.StatusNoContent,
		},
		{
			name:         "session not found",
			refreshToken: refreshToken,
			withCookie:   true,
			setupMock: func(sessionRepo *mocks.MockSessionRepository) {
				sessionRepo.EXPECT().
					FindByRefreshTokenHashForUpdate(mock.Anything, refreshTokenHash).
					Return(nil, errors.New(domain.ErrSessionNotFound)).
					Once()
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:         "session revoke failure",
			refreshToken: refreshToken,
			withCookie:   true,
			setupMock: func(sessionRepo *mocks.MockSessionRepository) {
				sessionRepo.EXPECT().
					FindByRefreshTokenHashForUpdate(mock.Anything, refreshTokenHash).
					Return(revokeFailureSession, nil).
					Once()

				sessionRepo.EXPECT().
					Revoke(mock.Anything, revokeFailureSession.ID, mock.Anything).
					Return(errors.New(domain.ErrDB)).
					Once()
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
}
