package tests

import (
	"context"
	"crypto/rand"
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

func TestDelete(t *testing.T) {
	t.Parallel()

	type chatRepositoryMockFunc func(mc *minimock.Controller) repository.ChatRepository

	type args struct {
		ctx context.Context
		id  int64
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

		id = gofakeit.Int64()

		serviceErr = fmt.Errorf("ошибка при тестировании")
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		err                error
		chatRepositoryMock chatRepositoryMockFunc
	}{
		{
			name: "Успешный вариант",
			args: args{
				ctx: ctx,
				id:  id,
			},
			err: nil,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(nil)
				return mock
			},
		},
		{
			name: "Неуспешный вариант",
			args: args{
				ctx: ctx,
				id:  id,
			},
			err: serviceErr,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repositoryMocks.NewChatRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepositoryMock := tt.chatRepositoryMock(mc)
			service := chat.NewChatService(userRepositoryMock)

			err := service.Delete(tt.args.ctx, tt.args.id)
			require.Equal(t, tt.err, err)
		})
	}
}
