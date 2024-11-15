package chat

import (
	"context"

	"github.com/vbulash/chat-server/internal/client/db"

	"github.com/vbulash/chat-server/internal/converter"
	"github.com/vbulash/chat-server/internal/model"
	"github.com/vbulash/chat-server/internal/repository"
	"github.com/vbulash/chat-server/internal/service"
)

type serviceLayer struct {
	repoLayer repository.ChatRepository
	txManager db.TxManager
}

// NewChatService Создание сервисного слоя
func NewChatService(repo repository.ChatRepository, txManager db.TxManager) service.ChatService {
	return &serviceLayer{
		repoLayer: repo,
		txManager: txManager,
	}
}

func (s *serviceLayer) CreateSend(ctx context.Context, info *model.ChatInfo) (int64, error) {
	var id int64
	var err error

	if nil == s.txManager { // Ветка для упрощенного юнит-теста
		id, err = s.repoLayer.CreateSend(ctx, converter.ModelChatInfoToDescChatInfo(info))
	} else {
		err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
			var err error
			id, err = s.repoLayer.CreateSend(ctx, converter.ModelChatInfoToDescChatInfo(info))
			if err != nil {
				return err
			}
			// ..
			return nil
		})
	}

	return id, err
}

func (s *serviceLayer) Get(ctx context.Context, id int64) (*model.Chat, error) {
	var chat *model.Chat
	var err error

	if nil == s.txManager { // Ветка для упрощенного юнит-теста
		chat, err = s.repoLayer.Get(ctx, id)
	} else {
		err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
			var err error
			chat, err = s.repoLayer.Get(ctx, id)
			if err != nil {
				return err
			}
			// ..
			return nil
		})
	}

	return chat, err
}

func (s *serviceLayer) Change(ctx context.Context, id int64, info *model.ChatInfo) error {
	var err error

	if nil == s.txManager { // Ветка для упрощенного юнит-теста
		err = s.repoLayer.Change(ctx, id, converter.ModelChatInfoToDescChatInfo(info))
	} else {
		err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
			errTx := s.repoLayer.Change(ctx, id, converter.ModelChatInfoToDescChatInfo(info))
			if errTx != nil {
				return errTx
			}
			// ..
			return nil
		})
	}

	return err
}

func (s *serviceLayer) Delete(ctx context.Context, id int64) error {
	var err error

	if nil == s.txManager { // Ветка для упрощенного юнит-теста
		err = s.repoLayer.Delete(ctx, id)
	} else {
		err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
			errTx := s.repoLayer.Delete(ctx, id)
			if errTx != nil {
				return errTx
			}
			// ..
			return nil
		})
	}

	return err
}
