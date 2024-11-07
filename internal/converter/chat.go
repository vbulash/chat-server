package converter

import (
	"database/sql"

	"github.com/vbulash/chat-server/internal/model"
	desc "github.com/vbulash/chat-server/pkg/chat_v2"
)

// ModelRecipientsToDescRecipients Преобразование из модели в GRPC
func ModelRecipientsToDescRecipients(modelRecipients []*model.UserIdentity) []*desc.UserIdentity {
	recipients := make([]*desc.UserIdentity, len(modelRecipients))
	for index, item := range modelRecipients {
		recipients[index] = &desc.UserIdentity{
			Id:    item.ID,
			Name:  item.Name,
			Email: item.Email,
		}
	}
	return recipients
}

// DescRecipientsToModelRecipients Преобразование из GRPC в модель
func DescRecipientsToModelRecipients(descRecipients []*desc.UserIdentity) []*model.UserIdentity {
	recipients := make([]*model.UserIdentity, len(descRecipients))
	for index, item := range descRecipients {
		recipients[index] = &model.UserIdentity{
			ID:    item.Id,
			Name:  item.Name,
			Email: item.Email,
		}
	}
	return recipients
}

// ModelChatInfoToDescChatInfo Преобразование из модели в GRPC
func ModelChatInfoToDescChatInfo(info *model.ChatInfo) *desc.ChatInfo {
	return &desc.ChatInfo{
		Recipients: ModelRecipientsToDescRecipients(info.Recipients),
		Text:       info.Body,
	}
}

// DescChatInfoToModelChatInfo Преобразование из GRPC в модель
func DescChatInfoToModelChatInfo(info *desc.ChatInfo) *model.ChatInfo {
	return &model.ChatInfo{
		Recipients: DescRecipientsToModelRecipients(info.Recipients),
		Body:       info.Text,
	}
}

// DescChatToModelChat Преобразование из GRPC в модель
func DescChatToModelChat(info *desc.Chat) *model.Chat {
	var updateAt sql.NullTime
	_ = updateAt.Scan(info.UpdatedAt.AsTime())

	translatedInfo := DescChatInfoToModelChatInfo(info.Info)

	return &model.Chat{
		ID:        info.Id,
		Info:      *translatedInfo,
		CreatedAt: info.CreatedAt.AsTime(),
		UpdatedAt: updateAt,
	}
}
