package out

import (
	chatsqlc "GoatChat/GoatChat/internal/chat/adapter/out/sqlc"
	"GoatChat/GoatChat/internal/chat/domain"
	dbtx "GoatChat/GoatChat/internal/shared/db_tx"
	"context"
	"errors"
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

func (r *ChatPgRepository) Save(ctx context.Context, cr domain.ChatRoom) error {
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

func (r *ChatPgRepository) Delete(ctx context.Context, cId uuid.UUID) error {
	q := r.queries(ctx)
	return q.DeleteChatRoomByID(ctx, chatsqlc.DeleteChatRoomByIDParams{ID: cId})
}

func (r *ChatPgRepository) Update(ctx context.Context, cr domain.ChatRoom) error {
	q := r.queries(ctx)
	err := q.UpdateChatRoomByID(ctx, chatsqlc.UpdateChatRoomByIDParams{
		RoomType:    string(cr.RoomType),
		Name:        &cr.RoomName.Value,
		Description: &cr.Description.Value,
		// LastMessageID: ^cr.LastMessageID, TODO(mgyoo) : Message 테이블 추가시 구현
		ID: cr.ID,
	})
	if err != nil {
		return errors.New("failed to update chatroom reason: " + err.Error())
	}

	return nil
}

func (r *ChatPgRepository) FindByChatroomId(ctx context.Context, cId uuid.UUID) (*domain.ChatRoom, error) {
	q := r.queries(ctx)

	cfd, err := q.GetChatRoomByID(ctx, chatsqlc.GetChatRoomByIDParams{ID: cId})
	if err != nil {
		return nil, errors.New("failed to get chat room by id: " + cId.String())
	}

	cr := toDomain(cfd)
	if cr == nil {
		return nil, errors.New("failed to get chat room by id, deleted chatroom")
	}

	return cr, nil
}

func (r *ChatPgRepository) ExistsByChatroomId(ctx context.Context, cId uuid.UUID) bool {
	_, err := r.FindByChatroomId(ctx, cId)

	return err == nil
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

func toDomain(cfd chatsqlc.ChatRoom) *domain.ChatRoom {
	if cfd.DeletedAt.Valid {
		return nil
	}

	return &domain.ChatRoom{
		ID:          cfd.ID,
		RoomType:    domain.RoomType(cfd.RoomType),
		RoomName:    domain.RoomName{Value: *cfd.Name},
		Description: domain.RoomDescription{Value: *cfd.Description},
		OwnerId:     domain.RoomOwnerId{Value: cfd.OwnerID.String()},
	}
}
