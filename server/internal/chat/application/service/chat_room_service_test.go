package service_test

import (
	"GoatChat/GoatChat/internal/chat/application/port/mocks"
	"GoatChat/GoatChat/internal/chat/application/service"
	domain "GoatChat/GoatChat/internal/chat/domain"
	"context"
	"testing"

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

	crsvc := service.NewChatService(m)

	ctx := context.TODO()
	cr, err := crsvc.CreateChatroom(ctx, domain.RoomTypeDirect, domain.RoomName{Value: "MGYOO"}, domain.RoomDescription{}, domain.RoomOwnerId{Value: "123142141"})
	require.NoError(t, err)

	assert.Equal(t, "MGYOO", cr.RoomName.Value)
}
