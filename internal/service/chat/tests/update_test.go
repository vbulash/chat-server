package tests

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"

	"github.com/vbulash/chat-server/internal/repository"

	"github.com/vbulash/chat-server/internal/converter"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/vbulash/chat-server/internal/model"
	repositoryMocks "github.com/vbulash/chat-server/internal/repository/mocks"
	"github.com/vbulash/chat-server/internal/service/chat"
	desc "github.com/vbulash/chat-server/pkg/chat_v2"
)

func TestUpdate(t *testing.T) {
	t.Parallel()

	type chatRepositoryMockFunc func(mc *minimock.Controller) repository.ChatRepository

	type args struct {
		ctx     context.Context
		id      int64
		request *model.ChatInfo
	}

	nBig, err := rand.Int(rand.Reader, big.NewInt(9))
	if err != nil {
		panic(err)
	}
	recipients := make([]*desc.UserIdentity, nBig.Int64()+1) // 1 .. 10
	for i := range recipients {
		recipients[i] = &desc.UserIdentity{
			Id:    int64(i),
			Name:  gofakeit.Name(),
			Email: gofakeit.Email(),
		}
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id   = gofakeit.Int64()
		text = gofakeit.Sentence(20)

		serviceErr = fmt.Errorf("ошибка при тестировании")

		request = &desc.ChatInfo{
			Recipients: recipients,
			Text:       text,
		}

		info = &model.ChatInfo{
			Recipients: converter.DescRecipientsToModelRecipients(recipients),
			Body:       text,
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		want               int64
		err                error
		chatRepositoryMock chatRepositoryMockFunc
	}{
		{
			name: "Успешный вариант",
			args: args{
				ctx:     ctx,
				id:      id,
				request: info,
			},
			want: id,
			err:  nil,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(mc)
				mock.ChangeMock.Expect(ctx, id, request).Return(nil)
				return mock
			},
		},
		{
			name: "Неуспешный вариант",
			args: args{
				ctx:     ctx,
				id:      id,
				request: info,
			},
			want: 0,
			err:  serviceErr,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(mc)
				mock.ChangeMock.Expect(ctx, id, request).Return(serviceErr)
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

			err := service.Change(tt.args.ctx, tt.args.id, tt.args.request)
			require.Equal(t, tt.err, err)
		})
	}
}
