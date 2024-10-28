package service

import (
	"context"
	"github.com/vbulash/chat-server/internal/model"
)

type ChatService interface {
	CreateSend(ctx context.Context, info *model.ChatInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.Chat, error)
	Change(ctx context.Context, id int64, info *model.ChatInfo) error
	Delete(ctx context.Context, id int64) error
}
