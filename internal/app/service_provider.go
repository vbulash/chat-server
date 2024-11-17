package app

import (
	"context"
	"log"

	"github.com/vbulash/platform_common/pkg/client/db"
	"github.com/vbulash/platform_common/pkg/client/db/pg"
	"github.com/vbulash/platform_common/pkg/client/db/transaction"
	"github.com/vbulash/platform_common/pkg/closer"

	api "github.com/vbulash/chat-server/internal/api/chat"
	chatAPI "github.com/vbulash/chat-server/internal/api/chat"
	"github.com/vbulash/chat-server/internal/repository"
	chatRepository "github.com/vbulash/chat-server/internal/repository/chat"
	"github.com/vbulash/chat-server/internal/service"
	chatService "github.com/vbulash/chat-server/internal/service/chat"
	"github.com/vbulash/platform_common/pkg/config"
	"github.com/vbulash/platform_common/pkg/config/env"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient     db.Client
	txManager    db.TxManager
	repoLayer    repository.ChatRepository
	serviceLayer service.ChatService
	apiLayer     *api.ChatsAPI
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// PGConfig ...
func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

// GRPCConfig ...
func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

// DBClient Клиент БД
func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		client, err := pg.New(ctx, s.PGConfig().DSN())
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

// TxManager Менеджер транзакций
func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

// RepoLayer Слой репозитория
func (s *serviceProvider) RepoLayer(ctx context.Context) repository.ChatRepository {
	if s.repoLayer == nil {
		repoLayer := chatRepository.NewChatRepository(s.DBClient(ctx))
		s.repoLayer = repoLayer
	}

	return s.repoLayer
}

// ServiceLayer Слой сервиса
func (s *serviceProvider) ServiceLayer(ctx context.Context) service.ChatService {
	if s.serviceLayer == nil {
		serviceLayer := chatService.NewChatService(
			s.RepoLayer(ctx),
			s.TxManager(ctx),
		)
		s.serviceLayer = serviceLayer
	}

	return s.serviceLayer
}

// APILayer Слой API
func (s *serviceProvider) APILayer(ctx context.Context) *api.ChatsAPI {
	if s.apiLayer == nil {
		apiLayer := chatAPI.NewAPI(s.ServiceLayer(ctx))
		s.apiLayer = apiLayer
	}

	return s.apiLayer
}
