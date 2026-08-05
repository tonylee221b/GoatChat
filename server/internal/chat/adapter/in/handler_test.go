package in_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	in "GoatChat/GoatChat/internal/chat/adapter/in"
	"GoatChat/GoatChat/internal/chat/application"
	"GoatChat/GoatChat/internal/chat/application/port/mocks"
	"GoatChat/GoatChat/internal/chat/domain"
	"GoatChat/GoatChat/internal/shared/db_tx/testutils"

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

			svc := application.NewChatService(testutils.StubTx{}, m)
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
					Save(mock.Anything, mock.Anything).
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
					Save(mock.Anything, mock.Anything).
					Return(errors.New("save error")).
					Once()
			},
			wantStatus: http.StatusBadRequest,
		},
	}
}

func TestDeleteChatroom(t *testing.T) {
	roomID := uuid.New()
	tests := []struct {
		name       string
		body       in.ChatroomDeleteRequest
		setupMock  func(*mocks.MockChatRepository)
		wantStatus int
	}{
		{
			name: "success",
			body: in.ChatroomDeleteRequest{RoomID: roomID.String()},
			setupMock: func(m *mocks.MockChatRepository) {
				m.EXPECT().
					Delete(mock.Anything, roomID).
					Return(nil).
					Once()
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid room id",
			body:       in.ChatroomDeleteRequest{RoomID: "invalid room id"},
			setupMock:  func(*mocks.MockChatRepository) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "service error",
			body: in.ChatroomDeleteRequest{RoomID: roomID.String()},
			setupMock: func(m *mocks.MockChatRepository) {
				m.EXPECT().
					Delete(mock.Anything, roomID).
					Return(errors.New("delete error")).
					Once()
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			require.NoError(t, json.NewEncoder(&buf).Encode(tt.body))

			req := httptest.NewRequest(http.MethodDelete, "/MOCK_URL", &buf)
			req.Header.Set("Content-Type", "application/json")

			m := mocks.NewMockChatRepository(t)
			tt.setupMock(m)
			svc := application.NewChatService(testutils.StubTx{}, m)
			h := in.NewChatHandler(*svc)

			rec := httptest.NewRecorder()
			h.DeleteChatroom(rec, req)

			res := rec.Result()
			defer res.Body.Close()
			require.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestUpdateChatroom(t *testing.T) {
	roomID := uuid.New()
	tests := []struct {
		name       string
		body       in.ChatroomUpdateRequest
		setupMock  func(*mocks.MockChatRepository)
		wantStatus int
	}{
		{
			name: "success",
			body: in.ChatroomUpdateRequest{
				RoomID:      roomID.String(),
				RoomType:    "group",
				RoomName:    "updated room",
				Description: "updated description",
			},
			setupMock: func(m *mocks.MockChatRepository) {
				m.EXPECT().
					FindByChatroomId(mock.Anything, roomID).
					Return(&domain.ChatRoom{ID: roomID}, nil).
					Once()
				m.EXPECT().
					Update(mock.Anything, mock.MatchedBy(func(room domain.ChatRoom) bool {
						return room.ID == roomID &&
							room.RoomType == domain.RoomTypeGroup &&
							room.RoomName.Value == "updated room" &&
							room.Description.Value == "updated description"
					})).
					Return(nil).
					Once()
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "invalid room id",
			body: in.ChatroomUpdateRequest{
				RoomID:   "invalid room id",
				RoomType: "group",
				RoomName: "updated room",
			},
			setupMock:  func(*mocks.MockChatRepository) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid room type",
			body: in.ChatroomUpdateRequest{
				RoomID:   roomID.String(),
				RoomType: "invalid",
				RoomName: "updated room",
			},
			setupMock:  func(*mocks.MockChatRepository) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid room name",
			body: in.ChatroomUpdateRequest{
				RoomID:   roomID.String(),
				RoomType: "group",
				RoomName: "",
			},
			setupMock:  func(*mocks.MockChatRepository) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid room description",
			body: in.ChatroomUpdateRequest{
				RoomID:      roomID.String(),
				RoomType:    "group",
				RoomName:    "updated room",
				Description: strings.Repeat("a", domain.DescriptionMaxSize+1),
			},
			setupMock:  func(*mocks.MockChatRepository) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "service error",
			body: in.ChatroomUpdateRequest{
				RoomID:   roomID.String(),
				RoomType: "group",
				RoomName: "updated room",
			},
			setupMock: func(m *mocks.MockChatRepository) {
				m.EXPECT().
					FindByChatroomId(mock.Anything, roomID).
					Return(nil, errors.New("find error")).
					Once()
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			require.NoError(t, json.NewEncoder(&buf).Encode(tt.body))

			req := httptest.NewRequest(http.MethodPatch, "/MOCK_URL", &buf)
			req.Header.Set("Content-Type", "application/json")

			m := mocks.NewMockChatRepository(t)
			tt.setupMock(m)
			svc := application.NewChatService(testutils.StubTx{}, m)
			h := in.NewChatHandler(*svc)

			rec := httptest.NewRecorder()
			h.UpdateChatroom(rec, req)

			res := rec.Result()
			defer res.Body.Close()
			require.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}
