package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/vbulash/chat-server/config"
	"github.com/vbulash/chat-server/database/operations"
	"log"
	"math/big"
	"time"

	"github.com/brianvoe/gofakeit"
	chat "github.com/vbulash/chat-server/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const address = "localhost:50052"

func closeConnection(conn *grpc.ClientConn) {
	err := conn.Close()
	if err != nil {
		log.Fatalf("Фатальная ошибка закрытия коннекта к серверу: %v", err)
	}
}

func main() {
	// Week 2
	config.Config = config.LoadConfig()

	db, err := operations.InitDb()
	if err != nil {
		log.Fatalf("Фатальная ошибка коннекта к базе данных: %v", err)
	}

	//operations.Seed(db)

	notes, err := operations.Get(db)
	if err != nil {
		log.Fatalf("Фатальная ошибка получения данных из базs данных: %v", err)
	}
	log.Println("Получены записи из таблицы notes:")
	for index := range *notes {
		log.Printf("%#v", (*notes)[index])
	}

	// Week 1
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Фатальная ошибка коннекта к серверу: %v", err)
	}
	defer closeConnection(conn)

	client := chat.NewChatV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Create
	fmt.Println("Клиент: создание нового чата")
	nBig, err := rand.Int(rand.Reader, big.NewInt(9))
	if err != nil {
		panic(err)
	}
	usernames := make([]string, nBig.Int64()+1) // 1 .. 10
	for i := range usernames {
		usernames[i] = gofakeit.Question()
	}
	newRecord := &chat.CreateRequest{
		Usernames: usernames,
	}
	fmt.Printf("Клиент: создаем новый чат: %+v\n", newRecord)
	response1, err := client.Create(ctx, newRecord)
	if err != nil {
		log.Fatalf("Клиент: фатальная ошибка создания нового чата: %v", err)
	}
	fmt.Printf("Клиент: создан новый чат ID = %d\n", response1.Id)
	id := response1.Id // Сквозной ID по всем эндпойнтам

	// Delete
	fmt.Println()
	fmt.Println("Клиент: удаление чата из системы")
	fmt.Printf("Клиент: удаляем чат по ID = %d\n", id)
	_, err = client.Delete(ctx, &chat.DeleteRequest{Id: id})
	if err != nil {
		log.Fatalf("Клиент: фатальная ошибка удаления чата ID = %d: %v", id, err)
	}
	fmt.Printf("Клиент: чат ID = %d удалён\n", id)

	// SendMessage
	fmt.Println()
	fmt.Println("Клиент: отправка сообщения на сервер")
	record := &chat.SendMessageRequest{
		From:      gofakeit.Name(),
		Text:      gofakeit.Question(),
		Timestamp: timestamppb.New(gofakeit.Date()),
	}
	fmt.Printf("Клиент: отправляем сообщение на сервер: %+v\n", record)
	_, err = client.SendMessage(ctx, record)
	if err != nil {
		log.Fatalf("Клиент: фатальная ошибка отправления сообщения на сервер: %v", err)
	}
	fmt.Println("Клиент: сообщение отправлено на сервер")

}
