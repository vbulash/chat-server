package repository

import (
	"context"

	outermodel "github.com/vbulash/chat-server/internal/model"
	desc "github.com/vbulash/chat-server/pkg/chat_v2"
)

// ChatRepository Репо чата
type ChatRepository interface {
	CreateSend(ctx context.Context, request *desc.ChatInfo) (int64, error)
	Get(ctx context.Context, id int64) (*outermodel.Chat, error)
	Change(ctx context.Context, id int64, request *desc.ChatInfo) error
	Delete(_ context.Context, id int64) error
}
