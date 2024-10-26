package chat

import (
	"context"
	"encoding/json"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vbulash/chat-server/internal/repository"
	"github.com/vbulash/chat-server/internal/repository/chat/model"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/vbulash/chat-server/pkg/chat_v2"
)

type repo struct {
	db *pgxpool.Pool
}

// NewChatRepository Создание репо
func NewChatRepository(db *pgxpool.Pool) repository.ChatRepository {
	return &repo{db: db}
}

func (r *repo) CreateSend(ctx context.Context, request *desc.ChatInfo) (int64, error) {
	recipients, err := json.Marshal(request.Recipients)
	if err != nil {
		return 0, nil
	}
	query, args, err := squirrel.Insert("chats").
		Columns("recipients", "body").
		Values(recipients, request.GetText()).
		Suffix("RETURNING \"id\"").
		ToSql()
	if err != nil {
		return 0, nil
	}

	var id int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, nil
	}
	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*desc.Chat, error) {
	query, args, err := squirrel.
		Select("id, recipients, body, created_at, updated_at").
		From("chats").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}
	var chat model.Chat
	err = r.db.QueryRow(ctx, query, args).Scan(&chat)
	if err != nil {
		return nil, err
	}

	// Inline converter
	var updatedAt *timestamppb.Timestamp
	if chat.UpdatedAt.Valid {
		updatedAt = timestamppb.New(chat.UpdatedAt.Time)
	}

	recipients := make([]*desc.UserIdentity, len(chat.Info.Recipients))
	for index, item := range chat.Info.Recipients {
		recipients[index] = &desc.UserIdentity{
			Id:    item.ID,
			Name:  item.Name,
			Email: item.Email,
		}
	}

	return &desc.Chat{
		Id: chat.ID,
		Info: &desc.ChatInfo{
			Recipients: recipients,
			Text:       chat.Info.Body,
		},
		CreatedAt: timestamppb.New(chat.CreatedAt),
		UpdatedAt: updatedAt,
	}, nil
}

func (r *repo) Change(ctx context.Context, id int64, request *desc.ChatInfo) error {
	updates := make(map[string]interface{})
	if len(request.Recipients) > 0 {
		recipients := make([]model.UserIdentity, len(request.Recipients))
		for index, item := range request.Recipients {
			recipients[index] = model.UserIdentity{
				ID:    item.Id,
				Name:  item.Name,
				Email: item.Email,
			}
		}
		updates["recipients"] = recipients
	}
	if len(request.Text) > 0 {
		updates["text"] = request.Text
	}

	query, args, err := squirrel.Update("chats").
		SetMap(updates).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, query, args...)
	return err
}

func (r *repo) Delete(_ context.Context, id int64) error {
	_, err := squirrel.Delete("chats").Where(squirrel.Eq{"id": id}).Exec()
	return err
}
