package converter

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/vbulash/chat-server/internal/model"
	modelRepo "github.com/vbulash/chat-server/internal/repository/chat/redis/model"
)

// ToChatFromRepo ...
func ToChatFromRepo(chat *modelRepo.Chat) *model.Chat {
	var updatedAt sql.NullTime
	if chat.UpdatedAt != nil {
		updatedAt = sql.NullTime{
			Time:  time.Unix(0, *chat.UpdatedAt),
			Valid: true,
		}
	}

	var users []model.UserIdentity
	err := json.Unmarshal([]byte(chat.Recipients), &users)
	if err != nil {
		return nil
	}
	recipients := make([]*model.UserIdentity, len(users))
	for index, item := range users {
		recipients[index] = &model.UserIdentity{
			ID:    item.ID,
			Name:  item.Name,
			Email: item.Email,
		}
	}

	return &model.Chat{
		ID: chat.ID,
		Info: model.ChatInfo{
			Recipients: recipients,
			Body:       chat.Body,
		},
		CreatedAt: time.Unix(0, chat.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

// ToChatFromService ...
func ToChatFromService(chat *model.Chat) *modelRepo.Chat {
	var updatedAt int64
	if chat.UpdatedAt.Valid {
		updatedAt = chat.UpdatedAt.Time.Unix()
	}

	recipients, err := json.Marshal(chat.Info.Recipients)
	if err != nil {
		return nil
	}

	return &modelRepo.Chat{
		ID:         chat.ID,
		Recipients: string(recipients),
		CreatedAt:  chat.CreatedAt.Unix(),
		UpdatedAt:  &updatedAt,
	}
}
