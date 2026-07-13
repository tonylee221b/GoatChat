package in_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	in "GoatChat/GoatChat/internal/chat/adapter/in"
	"GoatChat/GoatChat/internal/chat/application"
	"GoatChat/GoatChat/internal/chat/application/port/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateChatroom(t *testing.T) {
	for _, tt := range createChatroomTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			require.NoError(t, json.NewEncoder(&buf).Encode(tt.body))

			req := httptest.NewRequest(http.MethodPost, "/MOCK_URL", &buf)
			req.Header.Set("Content-Type", "application/json")

			m := mocks.NewMockChatRepository(t)
			if tt.name != "service_error" {
				tt.setupMock(m)
			}

			svc := application.NewChatService(m)
			h := in.NewChatHandler(*svc)

			rec := httptest.NewRecorder()
			h.CreateChatroom(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			require.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

type createChatroomTestCase struct {
	name       string
	body       in.ChatroomCreateRequest
	setupMock  func(*mocks.MockChatRepository)
	wantStatus int
}

func createChatroomTestCases() []createChatroomTestCase {
	oid := uuid.NewString()
	return []createChatroomTestCase{
		{
			name: "success",
			body: in.ChatroomCreateRequest{
				OwnerID:  oid,
				RoomType: "direct",
				RoomName: "test room",
			},
			setupMock: func(m *mocks.MockChatRepository) {
				m.EXPECT().
					SaveChatroom(mock.Anything, mock.Anything).
					Return(nil).
					Once()
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "invalid room type",
			body: in.ChatroomCreateRequest{
				OwnerID:  oid,
				RoomType: "wrong",
				RoomName: "test room",
			},
			setupMock:  func(m *mocks.MockChatRepository) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid room name",
			body: in.ChatroomCreateRequest{
				OwnerID:  oid,
				RoomType: "direct",
				RoomName: "",
			},
			setupMock:  func(m *mocks.MockChatRepository) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid owner id",
			body: in.ChatroomCreateRequest{
				OwnerID:  "invalid owner id",
				RoomType: "direct",
				RoomName: "test room",
			},
			setupMock:  func(m *mocks.MockChatRepository) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "service error",
			body: in.ChatroomCreateRequest{
				OwnerID:  oid,
				RoomType: "direct",
				RoomName: "test room",
			},
			setupMock: func(m *mocks.MockChatRepository) {
				m.EXPECT().
					SaveChatroom(mock.Anything, mock.Anything).
					Return(errors.New("save error")).
					Once()
			},
			wantStatus: http.StatusBadRequest,
		},
	}
}
