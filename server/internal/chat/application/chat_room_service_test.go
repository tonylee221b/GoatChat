package application_test

import (
	"context"
	"errors"
	"testing"

	"GoatChat/GoatChat/internal/chat/application"
	"GoatChat/GoatChat/internal/chat/application/port/mocks"
	domain "GoatChat/GoatChat/internal/chat/domain"
	"GoatChat/GoatChat/internal/shared/db_tx/testutils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateChatroom(t *testing.T) {
	roomType := domain.RoomTypeDirect
	roomName := domain.RoomName{Value: "test room"}
	description := domain.RoomDescription{Value: "test description"}
	ownerID := domain.RoomOwnerId{Value: uuid.NewString()}

	for _, tt := range createChatroomTestCases(roomType, roomName, description, ownerID) {
		t.Run(tt.name, func(t *testing.T) {
			m := mocks.NewMockChatRepository(t)
			tt.setupMock(m)

			svc := application.NewChatService(testutils.StubTx{}, m)
			room, err := svc.CreateChatroom(
				context.Background(),
				roomType,
				roomName,
				description,
				ownerID,
			)

			if tt.wantErr == "" {
				require.NoError(t, err)
				require.Equal(t, roomName, room.RoomName)
				return
			}
			require.Nil(t, room)
			require.EqualError(t, err, tt.wantErr)
		})
	}
}

func TestDeleteChatroom(t *testing.T) {
	roomID := uuid.New()

	for _, tt := range deleteChatroomTestCases(roomID) {
		t.Run(tt.name, func(t *testing.T) {
			m := mocks.NewMockChatRepository(t)
			tt.setupMock(m)

			svc := application.NewChatService(testutils.StubTx{}, m)
			err := svc.DeleteChatroom(context.Background(), roomID)

			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.EqualError(t, err, tt.wantErr)
		})
	}
}

func TestUpdateChatroom(t *testing.T) {
	roomID := uuid.New()
	updatedType := domain.RoomTypeGroup
	updatedName := domain.RoomName{Value: "updated room"}
	updatedDescription := domain.RoomDescription{Value: "updated description"}

	for _, tt := range updateChatroomTestCases(roomID, updatedType, updatedName, updatedDescription) {
		t.Run(tt.name, func(t *testing.T) {
			m := mocks.NewMockChatRepository(t)
			tt.setupMock(m)

			svc := application.NewChatService(testutils.StubTx{}, m)
			err := svc.UpdateChatroom(
				context.Background(),
				roomID,
				updatedType,
				updatedName,
				updatedDescription,
			)

			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.EqualError(t, err, tt.wantErr)
		})
	}
}

type chatroomServiceTestCase struct {
	name      string
	setupMock func(*mocks.MockChatRepository)
	wantErr   string
}

func createChatroomTestCases(
	roomType domain.RoomType,
	roomName domain.RoomName,
	description domain.RoomDescription,
	ownerID domain.RoomOwnerId,
) []chatroomServiceTestCase {
	return []chatroomServiceTestCase{
		{
			name: "success",
			setupMock: func(m *mocks.MockChatRepository) {
				m.EXPECT().
					Save(mock.Anything, mock.MatchedBy(func(room domain.ChatRoom) bool {
						return room.RoomType == roomType &&
							room.RoomName == roomName &&
							room.Description == description &&
							room.OwnerId == ownerID
					})).
					Return(nil).
					Once()
			},
		},
		{
			name: "repository error",
			setupMock: func(m *mocks.MockChatRepository) {
				m.EXPECT().
					Save(mock.Anything, mock.Anything).
					Return(errors.New("save error")).
					Once()
			},
			wantErr: "failed to save chat room, error: save error",
		},
	}
}

func deleteChatroomTestCases(roomID uuid.UUID) []chatroomServiceTestCase {
	return []chatroomServiceTestCase{
		{
			name: "success",
			setupMock: func(m *mocks.MockChatRepository) {
				m.EXPECT().
					Delete(mock.Anything, roomID).
					Return(nil).
					Once()
			},
		},
		{
			name: "repository error",
			setupMock: func(m *mocks.MockChatRepository) {
				m.EXPECT().
					Delete(mock.Anything, roomID).
					Return(errors.New("delete error")).
					Once()
			},
			wantErr: "delete error",
		},
	}
}

func updateChatroomTestCases(
	roomID uuid.UUID,
	updatedType domain.RoomType,
	updatedName domain.RoomName,
	updatedDescription domain.RoomDescription,
) []chatroomServiceTestCase {
	newChatroom := func() *domain.ChatRoom {
		return &domain.ChatRoom{
			ID:          roomID,
			RoomType:    domain.RoomTypeDirect,
			RoomName:    domain.RoomName{Value: "old room"},
			Description: domain.RoomDescription{Value: "old description"},
			OwnerId:     domain.RoomOwnerId{Value: uuid.NewString()},
		}
	}

	return []chatroomServiceTestCase{
		{
			name: "success",
			setupMock: func(m *mocks.MockChatRepository) {
				m.EXPECT().
					FindByChatroomId(mock.Anything, roomID).
					Return(newChatroom(), nil).
					Once()
				m.EXPECT().
					Update(mock.Anything, mock.MatchedBy(func(room domain.ChatRoom) bool {
						return room.ID == roomID &&
							room.RoomType == updatedType &&
							room.RoomName == updatedName &&
							room.Description == updatedDescription
					})).
					Return(nil).
					Once()
			},
		},
		{
			name: "chatroom lookup error",
			setupMock: func(m *mocks.MockChatRepository) {
				m.EXPECT().
					FindByChatroomId(mock.Anything, roomID).
					Return(nil, errors.New("find error")).
					Once()
			},
			wantErr: "failed to update chat room, error: find error",
		},
		{
			name: "repository update error",
			setupMock: func(m *mocks.MockChatRepository) {
				m.EXPECT().
					FindByChatroomId(mock.Anything, roomID).
					Return(newChatroom(), nil).
					Once()
				m.EXPECT().
					Update(mock.Anything, mock.Anything).
					Return(errors.New("update error")).
					Once()
			},
			wantErr: "failed to update chat room, error: update error",
		},
	}
}
