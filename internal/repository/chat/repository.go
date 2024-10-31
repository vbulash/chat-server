package chat

import (
	"context"
	"encoding/json"
	"time"

	"github.com/vbulash/chat-server/internal/client/db"

	"github.com/Masterminds/squirrel"
	"github.com/vbulash/chat-server/internal/repository"
	"github.com/vbulash/chat-server/internal/repository/chat/model"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/vbulash/chat-server/pkg/chat_v2"
)

type repoLayer struct {
	db db.Client
}

// NewChatRepository Создание репо
func NewChatRepository(db db.Client) repository.ChatRepository {
	return &repoLayer{db: db}
}

func (r *repoLayer) CreateSend(ctx context.Context, request *desc.ChatInfo) (int64, error) {
	recipients := make([]model.UserIdentity, len(request.Recipients))
	for index, item := range request.Recipients {
		recipients[index] = model.UserIdentity{
			ID:    item.Id,
			Name:  item.Name,
			Email: item.Email,
		}
	}
	jsonData, _ := json.Marshal(recipients)
	query, args, err := squirrel.Insert("chats").
		Columns("recipients", "body").
		Values(string(jsonData), request.GetText()).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return 0, nil
	}

	q := db.Query{
		Name:     "chat-server.CreateSend",
		QueryRaw: query,
	}
	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return 0, nil
	}
	return id, nil
}

func (r *repoLayer) Get(ctx context.Context, id int64) (*desc.Chat, error) {
	query, args, err := squirrel.
		Select("id, recipients, body, created_at, updated_at").
		From("chats").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "chat-server.Get",
		QueryRaw: query,
	}
	var chat model.Chat
	err = r.db.DB().QueryRowContext(ctx, q, args...).
		Scan(&chat.ID, &chat.Info.Recipients, &chat.Info.Body, &chat.CreatedAt, &chat.UpdatedAt)
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

func (r *repoLayer) Change(ctx context.Context, id int64, request *desc.ChatInfo) error {
	bUpdated := false
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
		jsonData, _ := json.Marshal(recipients)
		updates["recipients"] = string(jsonData)
		bUpdated = true
	}
	if len(request.Text) > 0 {
		updates["body"] = request.Text
		bUpdated = true
	}
	if bUpdated {
		updates["updated_at"] = time.Now()
	}

	query, args, err := squirrel.Update("chats").
		SetMap(updates).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "chat-server.Change",
		QueryRaw: query,
	}
	_, err = r.db.DB().ExecContext(ctx, q, args...)
	return err
}

func (r *repoLayer) Delete(ctx context.Context, id int64) error {
	query, args, err := squirrel.Delete("chats").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "chat-server.Delete",
		QueryRaw: query,
	}
	_, err = r.db.DB().ExecContext(ctx, q, args...)
	return err
}
