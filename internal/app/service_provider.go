package app

import (
	"context"
	"log"

	"github.com/vbulash/chat-server/internal/client/db"
	"github.com/vbulash/chat-server/internal/client/db/pg"
	"github.com/vbulash/chat-server/internal/closer"

	api "github.com/vbulash/chat-server/internal/api/chat"
	chatAPI "github.com/vbulash/chat-server/internal/api/chat"
	"github.com/vbulash/chat-server/internal/config"
	"github.com/vbulash/chat-server/internal/repository"
	chatRepository "github.com/vbulash/chat-server/internal/repository/chat"
	"github.com/vbulash/chat-server/internal/service"
	chatService "github.com/vbulash/chat-server/internal/service/chat"
)

type serviceProvider struct {
	env *config.Env

	dbClient     db.Client
	repoLayer    *repository.ChatRepository
	serviceLayer *service.ChatService
	apiLayer     *api.ChatsAPI
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// Env Конфигурация контейнера
func (s *serviceProvider) Env() *config.Env {
	if s.env == nil {
		env, err := config.LoadConfig()
		if err != nil {
			log.Fatalf("Ошибка загрузки .env: %v", err)
		}
		s.env = env
	}

	return s.env
}

// DBClient Клиент БД
func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		client, err := pg.New(ctx, s.Env().DSN)
		if err != nil {
			log.Fatalf("Ошибка создания db клиента: %v", err)
		}

		err = client.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("Ошибка пинга db: %v", err)
		}

		closer.Add(client.Close)

		s.dbClient = client
	}
	return s.dbClient
}

// RepoLayer Слой репозитория
func (s *serviceProvider) RepoLayer(ctx context.Context) *repository.ChatRepository {
	if s.repoLayer == nil {
		repoLayer := chatRepository.NewChatRepository(s.DBClient(ctx))
		s.repoLayer = &repoLayer
	}
	return s.repoLayer
}

// ServiceLayer Слой сервиса
func (s *serviceProvider) ServiceLayer(ctx context.Context) *service.ChatService {
	if s.serviceLayer == nil {
		serviceLayer := chatService.NewChatService(*s.RepoLayer(ctx))
		s.serviceLayer = &serviceLayer
	}
	return s.serviceLayer
}

// APILayer Слой API
func (s *serviceProvider) APILayer(ctx context.Context) *api.ChatsAPI {
	if s.apiLayer == nil {
		apiLayer := chatAPI.NewAPI(*s.ServiceLayer(ctx))
		s.apiLayer = apiLayer
	}
	return s.apiLayer
}
