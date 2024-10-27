package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/vbulash/chat-server/config"

	"github.com/brianvoe/gofakeit"
	chat "github.com/vbulash/chat-server/pkg/chat_v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func closeConnection(conn *grpc.ClientConn) {
	err := conn.Close()
	if err != nil {
		log.Fatalf("Фатальная ошибка закрытия коннекта к серверу: %v", err)
	}
}

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env: %v", err)
	}
	config.Config = conf

	address := fmt.Sprintf("%s:%d", config.Config.ServerHost, config.Config.ServerPort)
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Фатальная ошибка коннекта к серверу: %v", err)
	}
	defer closeConnection(conn)

	client := chat.NewChatV2Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// CreateSend
	fmt.Println("Клиент: создание и отправка чата")
	recipients := make([]*chat.UserIdentity, 0)
	for index := 0; index < 3; index++ {
		recipients = append(recipients, &chat.UserIdentity{
			Id:    gofakeit.Int64(),
			Name:  gofakeit.Name(),
			Email: gofakeit.Email(),
		})
	}
	record1 := &chat.CreateSendRequest{
		Recipients: recipients,
		Text:       gofakeit.Sentence(9),
	}
	fmt.Printf("Клиент: создаем новый чат: %+v\n", record1)
	response1, err := client.CreateSend(ctx, record1)
	if err != nil {
		log.Fatalf("Клиент: фатальная ошибка создания записи чата: %v", err)
	}
	fmt.Printf("Клиент: создан новый чат ID = %d\n", response1.Id)
	id := response1.Id // Сквозной ID по всем эндпойнтам

	// Get
	fmt.Println()
	fmt.Println("Клиент: получение чата")
	fmt.Printf("Клиент: получаем информацию чата ID = %d\n", id)
	response2, err := client.Get(ctx, &chat.GetRequest{Id: id})
	if err != nil {
		log.Fatalf("Клиент: фатальная ошибка получения записи чата ID = %d: %v", id, err)
	}
	fmt.Printf("Клиент: получен чат %+v\n", response2)

	// Change
	fmt.Println()
	fmt.Println("Клиент: изменение чата")
	recipients = make([]*chat.UserIdentity, 0)
	for index := 0; index < 2; index++ {
		recipients = append(recipients, &chat.UserIdentity{
			Id:    gofakeit.Int64(),
			Name:  gofakeit.Name(),
			Email: gofakeit.Email(),
		})
	}
	record2 := &chat.ChangeRequest{
		Id:         id,
		Recipients: recipients,
		Text:       gofakeit.Sentence(12),
	}
	fmt.Printf("Клиент: обновляем информацию чата ID = %d: %+v\n", id, record2)
	_, err = client.Change(ctx, record2)
	if err != nil {
		log.Fatalf("Клиент: фатальная ошибка обновления записи чата ID = %d: %v", id, err)
	}
	fmt.Printf("Клиент: обновлена запись чата ID = %d\n", id)

	// Delete
	fmt.Println()
	fmt.Println("Клиент: удаление чата")
	fmt.Printf("Клиент: удаляем запись чата ID = %d\n", id)
	_, err = client.Delete(ctx, &chat.DeleteRequest{Id: id})
	if err != nil {
		log.Fatalf("Клиент: фатальная ошибка удаления записи чата ID = %d: %v", id, err)
	}
	fmt.Printf("Клиент: запись чата ID = %d удалена\n", id)
}
