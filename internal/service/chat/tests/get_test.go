package tests

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"
	"testing"

	"github.com/vbulash/chat-server/internal/repository"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/vbulash/chat-server/internal/model"
	repositoryMocks "github.com/vbulash/chat-server/internal/repository/mocks"
	"github.com/vbulash/chat-server/internal/service/chat"
)

func TestGet(t *testing.T) {
	t.Parallel()

	type chatRepositoryMockFunc func(mc *minimock.Controller) repository.ChatRepository

	type args struct {
		ctx     context.Context
		id      int64
		request *model.Chat
	}

	nBig, err := rand.Int(rand.Reader, big.NewInt(9))
	if err != nil {
		panic(err)
	}
	recipients := make([]*model.UserIdentity, nBig.Int64()+1) // 1 .. 10
	for i := range recipients {
		recipients[i] = &model.UserIdentity{
			ID:    int64(i),
			Name:  gofakeit.Name(),
			Email: gofakeit.Email(),
		}
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = gofakeit.Int64()
		text      = gofakeit.Sentence(20)
		createdAt = gofakeit.Date()
		updatedAt = gofakeit.Date()

		serviceErr = fmt.Errorf("ошибка при тестировании")

		request = &model.Chat{
			ID: id,
			Info: model.ChatInfo{
				Recipients: recipients,
				Body:       text,
			},
			CreatedAt: createdAt,
			UpdatedAt: sql.NullTime{
				Valid: true,
				Time:  updatedAt,
			},
		}

		repoResponse    = request
		serviceResponse = repoResponse
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		want               *model.Chat
		err                error
		chatRepositoryMock chatRepositoryMockFunc
	}{
		{
			name: "Успешный вариант",
			args: args{
				ctx:     ctx,
				id:      id,
				request: request,
			},
			want: serviceResponse,
			err:  nil,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(repoResponse, nil)
				return mock
			},
		},
		{
			name: "Неуспешный вариант",
			args: args{
				ctx:     ctx,
				id:      id,
				request: request,
			},
			want: nil,
			err:  serviceErr,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepositoryMock := tt.chatRepositoryMock(mc)
			// Упрощенный вариант инициализации сервиса - без менеджера транзакций
			service := chat.NewChatService(userRepositoryMock, nil)

			response, err := service.Get(tt.args.ctx, tt.args.id)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, response)
		})
	}
}
