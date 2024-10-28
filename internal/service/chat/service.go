package chat

import (
	"context"

	"github.com/vbulash/chat-server/internal/converter"
	"github.com/vbulash/chat-server/internal/model"
	"github.com/vbulash/chat-server/internal/repository"
	"github.com/vbulash/chat-server/internal/service"
)

type serviceLayer struct {
	repoLayer repository.ChatRepository
}

// NewServiceLayer Создание сервисного слоя
func NewServiceLayer(repo repository.ChatRepository) service.ChatService {
	return &serviceLayer{repoLayer: repo}
}

func (s *serviceLayer) CreateSend(ctx context.Context, info *model.ChatInfo) (int64, error) {
	return s.repoLayer.CreateSend(ctx, converter.ModelChatInfoToDescChatInfo(info))
}

func (s *serviceLayer) Get(ctx context.Context, id int64) (*model.Chat, error) {
	nonConverted, err := s.repoLayer.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return converter.DescChatToModelChat(nonConverted), nil
}

func (s *serviceLayer) Change(ctx context.Context, id int64, info *model.ChatInfo) error {
	return s.repoLayer.Change(ctx, id, converter.ModelChatInfoToDescChatInfo(info))
}

func (s *serviceLayer) Delete(ctx context.Context, id int64) error {
	return s.repoLayer.Delete(ctx, id)
}
