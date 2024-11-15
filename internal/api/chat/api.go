package chat

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/vbulash/chat-server/internal/converter"
	"github.com/vbulash/chat-server/internal/model"
	"github.com/vbulash/chat-server/internal/service"
	desc "github.com/vbulash/chat-server/pkg/chat_v2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ChatsAPI Слой ChatsAPI
type ChatsAPI struct {
	desc.UnimplementedChatV2Server
	serviceLayer service.ChatService
}

// NewAPI Создание ChatsAPI
func NewAPI(serviceLayer service.ChatService) *ChatsAPI {
	return &ChatsAPI{serviceLayer: serviceLayer}
}

// CreateSend Создание чата
func (apiLayer *ChatsAPI) CreateSend(ctx context.Context, request *desc.CreateSendRequest) (*desc.CreateSendResponse, error) {
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
func (apiLayer *ChatsAPI) Get(ctx context.Context, request *desc.GetRequest) (*desc.GetResponse, error) {
	chatObj, err := apiLayer.serviceLayer.Get(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	var updatedAt *timestamppb.Timestamp
	if chatObj.UpdatedAt.Valid {
		updatedAt = timestamppb.New(chatObj.UpdatedAt.Time)
	}

	return &desc.GetResponse{
		Id:         chatObj.ID,
		Recipients: converter.ModelRecipientsToDescRecipients(chatObj.Info.Recipients),
		Text:       chatObj.Info.Body,
		CreatedAt:  timestamppb.New(chatObj.CreatedAt),
		UpdatedAt:  updatedAt,
	}, nil
}

// Change Изменение чата
func (apiLayer *ChatsAPI) Change(ctx context.Context, request *desc.ChangeRequest) (*empty.Empty, error) {
	err := apiLayer.serviceLayer.Change(ctx, request.Id, &model.ChatInfo{
		Recipients: converter.DescRecipientsToModelRecipients(request.Recipients),
		Body:       request.GetText(),
	})

	return &empty.Empty{}, err
}

// Delete Удаление чата
func (apiLayer *ChatsAPI) Delete(ctx context.Context, request *desc.DeleteRequest) (*empty.Empty, error) {
	err := apiLayer.serviceLayer.Delete(ctx, request.Id)

	return &empty.Empty{}, err
}
