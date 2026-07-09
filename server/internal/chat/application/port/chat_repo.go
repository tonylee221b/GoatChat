package port

import (
	"GoatChat/GoatChat/internal/chat/domain"
	"context"
)

type ChatRepository interface {
	SaveChatroom(ctx context.Context, cr domain.ChatRoom) error
}
