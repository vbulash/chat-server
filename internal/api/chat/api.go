package chat

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/vbulash/chat-server/internal/converter"
	"github.com/vbulash/chat-server/internal/model"
	"github.com/vbulash/chat-server/internal/service"
	desc "github.com/vbulash/chat-server/pkg/chat_v2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// API Слой API
type API struct {
	desc.UnimplementedChatV2Server
	serviceLayer service.ChatService
}

// NewAPI Создание API
func NewAPI(serviceLayer service.ChatService) *API {
	return &API{serviceLayer: serviceLayer}
}

// CreateSend Создание чата
func (apiLayer *API) CreateSend(ctx context.Context, request *desc.CreateSendRequest) (*desc.CreateSendResponse, error) {
	fmt.Println("Сервер: создание и отправка чата")

	id, err := apiLayer.serviceLayer.CreateSend(ctx, &model.ChatInfo{
		Recipients: converter.DescRecipientsToModelRecipients(request.Recipients),
		Body:       request.GetText(),
	})
	if err != nil {
		return nil, err
	}
	return &desc.CreateSendResponse{
		Id: id,
	}, nil
}

// Get Получение чата
func (apiLayer *API) Get(ctx context.Context, request *desc.GetRequest) (*desc.GetResponse, error) {
	fmt.Println("Сервер: получение чата")

	chatObj, err := apiLayer.serviceLayer.Get(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	var createdAt, updatedAt *timestamppb.Timestamp
	if chatObj.UpdatedAt.Valid {
		updatedAt = timestamppb.New(chatObj.UpdatedAt.Time)
	}
	createdAt = timestamppb.New(chatObj.CreatedAt)

	return &desc.GetResponse{
		Id:         chatObj.ID,
		Recipients: converter.ModelRecipientsToDescRecipients(chatObj.Info.Recipients),
		Text:       chatObj.Info.Body,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}, nil
}

// Change Изменение чата
func (apiLayer *API) Change(ctx context.Context, request *desc.ChangeRequest) (*empty.Empty, error) {
	fmt.Println("Сервер: обновление чата")

	err := apiLayer.serviceLayer.Change(ctx, request.Id, &model.ChatInfo{
		Recipients: converter.DescRecipientsToModelRecipients(request.Recipients),
		Body:       request.GetText(),
	})
	return &empty.Empty{}, err
}

// Delete Удаление чата
func (apiLayer *API) Delete(ctx context.Context, request *desc.DeleteRequest) (*empty.Empty, error) {
	fmt.Println("Сервер: удаление чата")

	err := apiLayer.serviceLayer.Delete(ctx, request.Id)
	return &empty.Empty{}, err
}
