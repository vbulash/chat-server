package app

import (
	"context"
	"log"

	"github.com/vbulash/platform_common/pkg/client/cache"

	redigo "github.com/gomodule/redigo/redis"
	chatRepositoryPg "github.com/vbulash/chat-server/internal/repository/chat/pg"
	chatRepositoryRedis "github.com/vbulash/chat-server/internal/repository/chat/redis"
	"github.com/vbulash/platform_common/pkg/client/cache/redis"

	"github.com/vbulash/platform_common/pkg/client/db"
	"github.com/vbulash/platform_common/pkg/client/db/pg"
	"github.com/vbulash/platform_common/pkg/client/db/transaction"
	"github.com/vbulash/platform_common/pkg/closer"

	api "github.com/vbulash/chat-server/internal/api/chat"
	chatAPI "github.com/vbulash/chat-server/internal/api/chat"
	"github.com/vbulash/chat-server/internal/repository"
	"github.com/vbulash/chat-server/internal/service"
	chatService "github.com/vbulash/chat-server/internal/service/chat"
	"github.com/vbulash/platform_common/pkg/config"
	"github.com/vbulash/platform_common/pkg/config/env"
)

type serviceProvider struct {
	pgConfig      config.PGConfig
	grpcConfig    config.GRPCConfig
	redisConfig   config.RedisConfig
	storageConfig config.StorageConfig

	redisPool   *redigo.Pool
	redisClient cache.RedisClient

	dbClient     db.Client
	txManager    db.TxManager
	repoLayer    repository.ChatRepository
	serviceLayer service.ChatService
	apiLayer     *api.ChatsAPI
}

const (
	redisMode = "redis"
	pgMode    = "pg"
)

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

func (s *serviceProvider) RedisConfig() config.RedisConfig {
	if s.redisConfig == nil {
		cfg, err := env.NewRedisConfig()
		if err != nil {
			log.Fatalf("failed to get redis config: %s", err.Error())
		}

		s.redisConfig = cfg
	}

	return s.redisConfig
}

func (s *serviceProvider) StorageConfig() config.StorageConfig {
	if s.storageConfig == nil {
		cfg, err := env.NewStorageConfig()
		if err != nil {
			log.Fatalf("failed to get storage config: %s", err.Error())
		}

		s.storageConfig = cfg
	}

	return s.storageConfig
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

func (s *serviceProvider) RedisPool() *redigo.Pool {
	if s.redisPool == nil {
		s.redisPool = &redigo.Pool{
			MaxIdle:     s.RedisConfig().MaxIdle(),
			IdleTimeout: s.RedisConfig().IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", s.RedisConfig().Address())
			},
		}
	}

	return s.redisPool
}

func (s *serviceProvider) RedisClient() cache.RedisClient {
	if s.redisClient == nil {
		s.redisClient = redis.NewClient(s.RedisPool(), s.RedisConfig())
	}

	return s.redisClient
}

// RepoLayer Слой репозитория
func (s *serviceProvider) RepoLayer(ctx context.Context) repository.ChatRepository {
	var repoLayer repository.ChatRepository
	if s.repoLayer == nil {
		switch s.StorageConfig().Mode() {
		case redisMode:
			repoLayer = chatRepositoryRedis.NewChatRepository(s.RedisClient())
			break
		case pgMode:
			repoLayer = chatRepositoryPg.NewChatRepository(s.DBClient(ctx))
			break
		default:
			repoLayer = nil
		}
		s.repoLayer = repoLayer
	}

	return s.repoLayer
}

// ServiceLayer Слой сервиса
func (s *serviceProvider) ServiceLayer(ctx context.Context) service.ChatService {
	var serviceLayer service.ChatService
	if s.serviceLayer == nil {
		switch s.StorageConfig().Mode() {
		case redisMode:
			serviceLayer = chatService.NewChatService(
				s.RepoLayer(ctx),
				nil,
			)
			break
		case pgMode:
			serviceLayer = chatService.NewChatService(
				s.RepoLayer(ctx),
				s.TxManager(ctx),
			)
			break
		default:
			serviceLayer = nil
		}
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
