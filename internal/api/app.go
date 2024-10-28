package api

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/vbulash/chat-server/internal/repository/chat"

	"github.com/jackc/pgx/v5/pgxpool"
	chat3 "github.com/vbulash/chat-server/internal/api/chat"
	chat2 "github.com/vbulash/chat-server/internal/service/chat"

	"github.com/vbulash/chat-server/config"

	desc "github.com/vbulash/chat-server/pkg/chat_v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// AppRun Инициализация и запуск приложения
func AppRun() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute) // Достаточно для интерактивной отладки
	defer cancel()

	conf, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("Ошибка загрузки .env: %v", err)
	}
	config.Config = conf

	poolConfig, err := pgxpool.ParseConfig(os.Getenv("DB_DSN"))
	if err != nil {
		return fmt.Errorf("Ошибка конфигурации pgxpool: %v", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return fmt.Errorf("Ошибка коннекта к БД: %v", err)
	}

	chatRepo := chat.NewChatRepository(pool)
	serviceLayer := chat2.NewServiceLayer(chatRepo)
	apiLayer := chat3.NewAPI(serviceLayer)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Config.ServerPort))
	if err != nil {
		return fmt.Errorf("Фатальная ошибка запуска / прослушивания: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterChatV2Server(s, apiLayer)

	log.Printf("Сервер прослушивает: %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		return fmt.Errorf("Фатальная ошибка запуска: %v", err)
	}

	return nil
}
