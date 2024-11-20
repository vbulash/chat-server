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

// NewChatService Создание сервисного слоя
func NewChatService(repo repository.ChatRepository) service.ChatService {
	return &serviceLayer{
		repoLayer: repo,
	}
}

func (s *serviceLayer) CreateSend(ctx context.Context, info *model.ChatInfo) (int64, error) {
	id, err := s.repoLayer.CreateSend(ctx, converter.ModelChatInfoToDescChatInfo(info))
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *serviceLayer) Get(ctx context.Context, id int64) (*model.Chat, error) {
	chat, err := s.repoLayer.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *serviceLayer) Change(ctx context.Context, id int64, info *model.ChatInfo) error {
	err := s.repoLayer.Change(ctx, id, converter.ModelChatInfoToDescChatInfo(info))
	if err != nil {
		return err
	}

	return nil
}

func (s *serviceLayer) Delete(ctx context.Context, id int64) error {
	err := s.repoLayer.Delete(ctx, id)
	if err != nil {
		return err
	}

	return err
}
