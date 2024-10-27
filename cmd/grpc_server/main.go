package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/vbulash/chat-server/internal/repository/chat"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vbulash/chat-server/internal/repository"

	"github.com/vbulash/chat-server/config"

	"github.com/golang/protobuf/ptypes/empty"
	desc "github.com/vbulash/chat-server/pkg/chat_v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	desc.UnimplementedChatV2Server
	chatRepository repository.ChatRepository
}

func (s *server) CreateSend(ctx context.Context, request *desc.CreateSendRequest) (*desc.CreateSendResponse, error) {
	fmt.Println("Сервер: создание и отправка чата")

	id, err := s.chatRepository.CreateSend(ctx, &desc.ChatInfo{
		Recipients: request.Recipients,
		Text:       request.Text,
	})
	if err != nil {
		return nil, err
	}
	return &desc.CreateSendResponse{
		Id: id,
	}, nil
}

func (s *server) Get(ctx context.Context, request *desc.GetRequest) (*desc.GetResponse, error) {
	fmt.Println("Сервер: получение чата")

	chatObj, err := s.chatRepository.Get(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return &desc.GetResponse{
		Id:         chatObj.Id,
		Recipients: chatObj.Info.Recipients,
		Text:       chatObj.Info.Text,
		CreatedAt:  chatObj.CreatedAt,
		UpdatedAt:  chatObj.UpdatedAt,
	}, nil
}

func (s *server) Change(ctx context.Context, request *desc.ChangeRequest) (*empty.Empty, error) {
	fmt.Println("Сервер: обновление чата")

	err := s.chatRepository.Change(ctx, request.Id, &desc.ChatInfo{
		Recipients: request.Recipients,
		Text:       request.Text,
	})
	return &empty.Empty{}, err
}

func (s *server) Delete(ctx context.Context, request *desc.DeleteRequest) (*empty.Empty, error) {
	fmt.Println("Сервер: удаление чата")

	err := s.chatRepository.Delete(ctx, request.Id)
	return &empty.Empty{}, err
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute) // Достаточно для интерактивной отладки
	defer cancel()

	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env: %v", err)
	}
	config.Config = conf

	poolConfig, err := pgxpool.ParseConfig(os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatalf("Ошибка конфигурации pgxpool: %v", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalf("Ошибка коннекта к БД: %v", err)
	}

	chatRepo := chat.NewChatRepository(pool)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Config.ServerPort))
	if err != nil {
		log.Fatalf("Фатальная ошибка запуска / прослушивания: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterChatV2Server(s, &server{
		chatRepository: chatRepo,
	})

	log.Printf("Сервер прослушивает: %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("Фатальная ошибка запуска: %v", err)
	}
}
