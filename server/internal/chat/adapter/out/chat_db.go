package out

import (
	chatsqlc "GoatChat/GoatChat/internal/chat/adapter/out/sqlc"
	"GoatChat/GoatChat/internal/chat/domain"
	dbtx "GoatChat/GoatChat/internal/shared/db_tx"
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatPgRepository struct {
	pool *pgxpool.Pool
}

func NewChatPgRepository(pool *pgxpool.Pool) *ChatPgRepository {
	return &ChatPgRepository{pool: pool}
}

func (r *ChatPgRepository) SaveChatroom(ctx context.Context, cr domain.ChatRoom) error {
	oid, err := uuid.Parse(cr.OwnerId.Value)
	if err != nil {
		slog.Error("failed to uuid parse, ", "error", err.Error())
		return err
	}

	q := r.queries(ctx)
	_, err = q.CreateChatRoom(ctx, chatsqlc.CreateChatRoomParams{
		ID:          cr.ID,
		RoomType:    string(cr.RoomType),
		Name:        &cr.RoomName.Value,
		Description: &cr.Description.Value,
		OwnerID:     oid,
	})

	return err
}

func (r *ChatPgRepository) queries(ctx context.Context) *chatsqlc.Queries {
	return chatsqlc.New(r.db(ctx))
}

func (r *ChatPgRepository) db(ctx context.Context) chatsqlc.DBTX {
	if pgxTx := dbtx.FromContext(ctx); pgxTx != nil {
		return pgxTx
	}

	return r.pool
}
