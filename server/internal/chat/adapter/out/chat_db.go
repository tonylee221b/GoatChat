package out

import (
	chatsqlc "GoatChat/GoatChat/internal/chat/adapter/out/sqlc"
	"GoatChat/GoatChat/internal/chat/domain"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type ChatPgRepository struct {
	q chatsqlc.Queries
}

func NewChatPgRepository(db chatsqlc.DBTX) *ChatPgRepository {
	return &ChatPgRepository{
		q: *chatsqlc.New(db),
	}
}

func (r *ChatPgRepository) SaveChatroom(ctx context.Context, cr domain.ChatRoom) error {
	oid, err := uuid.Parse(cr.OwnerId.Value)
	if err != nil {
		slog.Error("failed to uuid parse, ", "error", err.Error())
		return err
	}

	_, err = r.q.CreateChatRoom(ctx, chatsqlc.CreateChatRoomParams{
		ID:          cr.ID,
		RoomType:    string(cr.RoomType),
		Name:        &cr.RoomName.Value,
		Description: &cr.Description.Value,
		OwnerID:     oid,
	})

	return err
}
