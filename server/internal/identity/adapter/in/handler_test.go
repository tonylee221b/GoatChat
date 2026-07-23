package in_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

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

			m := mocks.NewMockUserRepository(t)
			if tt.setupMock != nil {
				tt.setupMock(m)
			}

			svc := application.NewUserService(testutils.StubTx{}, m)
			h := in.NewIdentityHandler(svc)

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
			m := mocks.NewMockUserRepository(t)
			if tt.setupMock != nil {
				tt.setupMock(m)
			}

			svc := application.NewUserService(testutils.StubTx{}, m)
			h := in.NewIdentityHandler(svc)

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
			m := mocks.NewMockUserRepository(t)
			if tt.setupMock != nil {
				tt.setupMock(m)
			}

			svc := application.NewUserService(testutils.StubTx{}, m)
			h := in.NewIdentityHandler(svc)

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

type registerUserTestCase struct {
	name       string
	body       in.RegisterUserReq
	setupMock  func(*mocks.MockUserRepository)
	wantStatus int
}

func registerUserTestCases() []registerUserTestCase {
	return []registerUserTestCase{
		{
			name: "success",
			body: in.RegisterUserReq{
				Username: "test-user",
			},
			setupMock: func(m *mocks.MockUserRepository) {
				m.EXPECT().
					ExistsByUsername(mock.Anything, mock.Anything).
					Return(false, nil)

				m.EXPECT().
					Save(mock.Anything, mock.Anything).
					Return(nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "invalid username",
			body: in.RegisterUserReq{
				Username: "",
			},
			setupMock:  nil,
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
	tu := "test-user"
	un, _ := domain.NewUsername(tu)
	u := domain.NewUser(un)

	return []findUserTestCase{
		{
			name:  "success",
			param: tu,
			setupMock: func(m *mocks.MockUserRepository) {
				m.EXPECT().
					FindByUsername(mock.Anything, mock.Anything).
					Return(u, nil).
					Once()
			},
			wantStatus: http.StatusOK,
		},
		{
			name:  "user not found",
			param: "wrong-user",
			setupMock: func(m *mocks.MockUserRepository) {
				m.EXPECT().
					FindByUsername(mock.Anything, mock.Anything).
					Return(nil, errors.New("user not found")).
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
	user := domain.NewUser(domainUsername)
	userNotFound := errors.New(domain.ErrUserNotFound)
	dbErr := errors.New(domain.ErrDB)

	validContact := in.UpdateUserContactReq{
		PhoneNumber: "010-1234-5678",
		Email:       "test@example.com",
	}

	return []updateContactTestCase{
		{
			name:  "success",
			param: username,
			body:  validContact,
			setupMock: func(m *mocks.MockUserRepository) {
				m.EXPECT().
					FindByUsername(mock.Anything, domainUsername).
					Return(user, nil).
					Once()

				m.EXPECT().
					Update(mock.Anything, mock.Anything).
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
			setupMock: func(m *mocks.MockUserRepository) {
				m.EXPECT().
					FindByUsername(mock.Anything, domainUsername).
					Return(nil, userNotFound).
					Once()
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:  "repository update error",
			param: username,
			body:  validContact,
			setupMock: func(m *mocks.MockUserRepository) {
				m.EXPECT().
					FindByUsername(mock.Anything, domainUsername).
					Return(user, nil).
					Once()

				m.EXPECT().
					Update(mock.Anything, mock.Anything).
					Return(dbErr).
					Once()
			},
			wantStatus: http.StatusBadRequest,
		},
	}
}
