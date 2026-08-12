package in

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"GoatChat/GoatChat/internal/chat/application"
	"GoatChat/GoatChat/internal/chat/application/port/mocks"
	"GoatChat/GoatChat/internal/shared/db_tx/testutils"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestChatRouterRoutesDeleteToDeleteChatroom(t *testing.T) {
	roomID := uuid.New()
	m := mocks.NewMockChatRepository(t)
	m.EXPECT().
		Delete(mock.Anything, roomID).
		Return(nil).
		Once()

	svc := application.NewChatService(testutils.StubTx{}, m)
	h := NewChatHandler(*svc)
	router := NewChatRouter(h)
	r := chi.NewRouter()
	router.route(r)

	var body bytes.Buffer
	require.NoError(t, json.NewEncoder(&body).Encode(ChatroomDeleteRequest{
		RoomID: roomID.String(),
	}))
	req := httptest.NewRequest(http.MethodDelete, ChatRouteGroup+ChatroomRouteGroup, &body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}
