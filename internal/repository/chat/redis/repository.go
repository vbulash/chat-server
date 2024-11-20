package redis

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	desc "github.com/vbulash/chat-server/pkg/chat_v2"

	redigo "github.com/gomodule/redigo/redis"

	"github.com/vbulash/chat-server/internal/model"
	"github.com/vbulash/chat-server/internal/repository"
	"github.com/vbulash/chat-server/internal/repository/chat/redis/converter"
	modelRepo "github.com/vbulash/chat-server/internal/repository/chat/redis/model"
	"github.com/vbulash/platform_common/pkg/client/cache"
)

type repoLayer struct {
	cl cache.RedisClient
}

// NewChatRepository Создание репо
func NewChatRepository(cl cache.RedisClient) repository.ChatRepository {
	return &repoLayer{cl: cl}
}

func (r *repoLayer) CreateSend(ctx context.Context, request *desc.ChatInfo) (int64, error) {
	recipients, err := json.Marshal(request.Recipients)
	if err != nil {
		return 0, nil
	}
	id := gofakeit.Int64()
	chat := &modelRepo.Chat{
		ID:         id,
		Recipients: string(recipients),
		Body:       request.Text,
		CreatedAt:  time.Now().UnixNano(),
	}

	idStr := strconv.FormatInt(id, 10)
	err = r.cl.HashSet(ctx, idStr, chat)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repoLayer) Get(ctx context.Context, id int64) (*model.Chat, error) {
	idStr := strconv.FormatInt(id, 10)
	values, err := r.cl.HGetAll(ctx, idStr)
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, model.ErrorChatNotFound
	}

	var chat modelRepo.Chat
	err = redigo.ScanStruct(values, &chat)
	if err != nil {
		return nil, err
	}

	return converter.ToChatFromRepo(&chat), nil
}

func (r *repoLayer) Change(ctx context.Context, id int64, request *desc.ChatInfo) error {
	// Нужно предварительно вычитать предыдущий CreatedAt
	stored, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	createdAt := stored.CreatedAt.UnixNano()
	recipients, err := json.Marshal(request.Recipients)
	if err != nil {
		return nil
	}
	idStr := strconv.FormatInt(id, 10)
	updatedAt := time.Now().UnixNano()

	chat := &modelRepo.Chat{
		ID:         id,
		Recipients: string(recipients),
		Body:       request.Text,
		CreatedAt:  createdAt,
		UpdatedAt:  &updatedAt,
	}

	err = r.cl.HashSet(ctx, idStr, chat)
	if err != nil {
		return err
	}

	return nil
}

func (r *repoLayer) Delete(ctx context.Context, id int64) error {
	idStr := strconv.FormatInt(id, 10)
	err := r.cl.Expire(ctx, idStr, 1)
	if err != nil {
		return err
	}

	return nil
}
