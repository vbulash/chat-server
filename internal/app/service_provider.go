package app

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	api "github.com/vbulash/chat-server/internal/api/chat"
	chatAPI "github.com/vbulash/chat-server/internal/api/chat"
	"github.com/vbulash/chat-server/internal/config"
	"github.com/vbulash/chat-server/internal/repository"
	chatRepository "github.com/vbulash/chat-server/internal/repository/chat"
	"github.com/vbulash/chat-server/internal/service"
	chatService "github.com/vbulash/chat-server/internal/service/chat"
)

type serviceProvider struct {
	env          *config.Env
	pool         *pgxpool.Pool
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

// Pool Пул соединений Postgres
func (s *serviceProvider) Pool(ctx context.Context) *pgxpool.Pool {
	if s.pool == nil {
		poolConfig, err := pgxpool.ParseConfig(s.Env().DSN)
		if err != nil {
			log.Fatalf("Ошибка конфигурации pgxpool: %v", err)
		}
		pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
		if err != nil {
			log.Fatalf("Ошибка коннекта к БД: %v", err)
		}

		s.pool = pool
	}
	return s.pool
}

// RepoLayer Слой репозитория
func (s *serviceProvider) RepoLayer(ctx context.Context) *repository.ChatRepository {
	if s.repoLayer == nil {
		repoLayer := chatRepository.NewChatRepository(s.Pool(ctx))
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
