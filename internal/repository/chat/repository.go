package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/vbulash/chat-server/internal/client/db"

	"github.com/Masterminds/squirrel"
	outermodel "github.com/vbulash/chat-server/internal/model"
	"github.com/vbulash/chat-server/internal/repository"
	innermodel "github.com/vbulash/chat-server/internal/repository/chat/model"
	desc "github.com/vbulash/chat-server/pkg/chat_v2"
)

const (
	tableName = "chats"

	idColumn         string = "id"
	recipientsColumn string = "recipients"
	bodyColumn       string = "body"
	createdAtColumn  string = "created_at"
	updatedAtColumn  string = "updated_at"
)

type repoLayer struct {
	db db.Client
}

// NewChatRepository Создание репо
func NewChatRepository(db db.Client) repository.ChatRepository {
	return &repoLayer{db: db}
}

func (r *repoLayer) CreateSend(ctx context.Context, request *desc.ChatInfo) (int64, error) {
	recipients := make([]innermodel.UserIdentity, len(request.Recipients))
	for index, item := range request.Recipients {
		recipients[index] = innermodel.UserIdentity{
			ID:    item.Id,
			Name:  item.Name,
			Email: item.Email,
		}
	}
	jsonData, _ := json.Marshal(recipients)
	query, args, err := squirrel.Insert(tableName).
		Columns(recipientsColumn, bodyColumn).
		Values(string(jsonData), request.GetText()).
		Suffix(fmt.Sprintf("RETURNING \"%s\"", idColumn)).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return 0, nil
	}

	q := db.Query{
		Name:     tableName + ".CreateSend",
		QueryRaw: query,
	}
	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return 0, nil
	}

	return id, nil
}

func (r *repoLayer) Get(ctx context.Context, id int64) (*outermodel.Chat, error) {
	query, args, err := squirrel.
		Select(strings.Join([]string{idColumn, recipientsColumn, bodyColumn, createdAtColumn, updatedAtColumn}, ", ")).
		From(tableName).
		Where(squirrel.Eq{idColumn: id}).
		PlaceholderFormat(squirrel.Dollar).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     tableName + ".Get",
		QueryRaw: query,
	}
	var chat innermodel.Chat
	err = r.db.DB().QueryRowContext(ctx, q, args...).
		// Разъяснение то же, что и в chat:
		// squirrel.Select("*") -> Scan(&chat)
		// squirrel.Select("id, recipients, body, created_at, updated_at") => Scan(&chat.ID, &chat.Info.Recipients, &chat.Info.Body, &chat.CreatedAt, &chat.UpdatedAt)
		Scan(&chat.ID, &chat.Info.Recipients, &chat.Info.Body, &chat.CreatedAt, &chat.UpdatedAt)
	if err != nil {
		return nil, err
	}

	recipients := make([]*outermodel.UserIdentity, len(chat.Info.Recipients))
	for index, item := range chat.Info.Recipients {
		recipients[index] = &outermodel.UserIdentity{
			ID:    item.ID,
			Name:  item.Name,
			Email: item.Email,
		}
	}

	// Преобразование внутренней модели во внешнюю - не стал выносить в конвертер
	// innermodel -> outermodel
	return &outermodel.Chat{
		ID: chat.ID,
		Info: outermodel.ChatInfo{
			Recipients: recipients,
			Body:       chat.Info.Body,
		},
		CreatedAt: chat.CreatedAt,
		UpdatedAt: chat.UpdatedAt,
	}, nil
}

func (r *repoLayer) Change(ctx context.Context, id int64, request *desc.ChatInfo) error {
	// Не стал выносить анализ заполнения полей - есть логика копирования получателей
	bUpdated := false
	updates := make(map[string]interface{})
	if len(request.Recipients) > 0 {
		recipients := make([]innermodel.UserIdentity, len(request.Recipients))
		for index, item := range request.Recipients {
			recipients[index] = innermodel.UserIdentity{
				ID:    item.Id,
				Name:  item.Name,
				Email: item.Email,
			}
		}
		jsonData, _ := json.Marshal(recipients)
		updates[recipientsColumn] = string(jsonData)
		bUpdated = true
	}
	if len(request.Text) > 0 {
		updates[bodyColumn] = request.Text
		bUpdated = true
	}
	if bUpdated {
		updates[updatedAtColumn] = time.Now()
	}

	query, args, err := squirrel.Update(tableName).
		SetMap(updates).
		Where(squirrel.Eq{idColumn: id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     tableName + ".Change",
		QueryRaw: query,
	}
	_, err = r.db.DB().ExecContext(ctx, q, args...)

	return err
}

func (r *repoLayer) Delete(ctx context.Context, id int64) error {
	query, args, err := squirrel.Delete(tableName).
		Where(squirrel.Eq{idColumn: id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     tableName + ".Delete",
		QueryRaw: query,
	}
	_, err = r.db.DB().ExecContext(ctx, q, args...)

	return err
}
