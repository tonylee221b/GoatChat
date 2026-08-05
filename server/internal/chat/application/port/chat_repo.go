package port

import (
	"GoatChat/GoatChat/internal/chat/domain"
	"context"

	"github.com/google/uuid"
)

type ChatRepository interface {
	Save(ctx context.Context, cr domain.ChatRoom) error
	Delete(ctx context.Context, cId uuid.UUID) error
	Update(ctx context.Context, cr domain.ChatRoom) error
	FindByChatroomId(ctx context.Context, cId uuid.UUID) (*domain.ChatRoom, error)
	ExistsByChatroomId(ctx context.Context, cId uuid.UUID) bool
}
