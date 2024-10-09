package main

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/golang/protobuf/ptypes/empty"
	chat "github.com/vbulash/chat-server/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const grpcPort = ":50052"

type server struct {
	chat.UnimplementedChatV1Server
}

func (s *server) Create(_ context.Context, _ *chat.CreateRequest) (*chat.CreateResponse, error) {
	fmt.Println("Сервер: создание нового чата")
	return &chat.CreateResponse{
		Id: gofakeit.Int64(),
	}, nil
}

func (s *server) Delete(_ context.Context, _ *chat.DeleteRequest) (*empty.Empty, error) {
	fmt.Println("Сервер: удаление чата из системы")
	return &empty.Empty{}, nil
}

func (s *server) SendMessage(_ context.Context, _ *chat.SendMessageRequest) (*empty.Empty, error) {
	fmt.Println("Сервер: отправка сообщения на сервер")
	return &empty.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s", grpcPort))
	if err != nil {
		log.Fatalf("Фатальная ошибка запуска / прослушивания: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	chat.RegisterChatV1Server(s, &server{})

	log.Printf("Сервер прослушивает: %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("Фатальная ошибка запуска: %v", err)
	}
}
