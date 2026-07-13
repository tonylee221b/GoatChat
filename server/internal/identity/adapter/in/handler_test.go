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

type registerUserTestCase struct {
	name       string
	body       in.UserRegisterRequest
	setupMock  func(*mocks.MockUserRepository)
	wantStatus int
}

func registerUserTestCases() []registerUserTestCase {
	return []registerUserTestCase{
		{
			name: "success",
			body: in.UserRegisterRequest{
				Username: "test-user",
			},
			setupMock: func(m *mocks.MockUserRepository) {
				m.EXPECT().
					Save(mock.Anything, mock.Anything).
					Return(nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "invalid username",
			body: in.UserRegisterRequest{
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
	u, _ := domain.NewUser(un)

	return []findUserTestCase{
		{
			name:  "success",
			param: tu,
			setupMock: func(m *mocks.MockUserRepository) {
				m.EXPECT().
					FindByUsername(mock.Anything, mock.Anything).
					Return(u, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:  "user not found",
			param: "wrong-user",
			setupMock: func(m *mocks.MockUserRepository) {
				m.EXPECT().
					FindByUsername(mock.Anything, mock.Anything).
					Return(nil, errors.New("user not found"))
			},
			wantStatus: http.StatusNotFound,
		},
	}
}
