package out

import (
	"testing"
	"time"

	chatsqlc "GoatChat/GoatChat/internal/chat/adapter/out/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestToDomainReturnsActiveChatroom(t *testing.T) {
	name := "test room"
	description := "test description"
	room := chatsqlc.ChatRoom{
		ID:          uuid.New(),
		RoomType:    "group",
		Name:        &name,
		Description: &description,
		OwnerID:     uuid.New(),
		DeletedAt:   pgtype.Timestamptz{},
	}

	got := toChatRoomDomain(room)

	require.NotNil(t, got)
	require.Equal(t, room.ID, got.ID)
}

func TestToDomainRejectsDeletedChatroom(t *testing.T) {
	name := "test room"
	description := "test description"
	room := chatsqlc.ChatRoom{
		ID:          uuid.New(),
		RoomType:    "group",
		Name:        &name,
		Description: &description,
		OwnerID:     uuid.New(),
		DeletedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	}

	got := toChatRoomDomain(room)

	require.Nil(t, got)
}
