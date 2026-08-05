package application_test

import (
	"context"
	"testing"

	"GoatChat/GoatChat/internal/chat/application"
	"GoatChat/GoatChat/internal/chat/application/port/mocks"
	domain "GoatChat/GoatChat/internal/chat/domain"
	"GoatChat/GoatChat/internal/shared/db_tx/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestChatroomCreate(t *testing.T) {
	m := mocks.NewMockChatRepository(t)
	m.EXPECT().
		SaveChatroom(mock.Anything, mock.Anything).
		Return(nil).
		Once()

	crsvc := application.NewChatService(testutils.StubTx{}, m)

	ctx := context.Background()
	cr, err := crsvc.CreateChatroom(ctx, domain.RoomTypeDirect, domain.RoomName{Value: "MGYOO"}, domain.RoomDescription{}, domain.RoomOwnerId{Value: "123142141"})
	require.NoError(t, err)

	assert.Equal(t, "MGYOO", cr.RoomName.Value)
}
