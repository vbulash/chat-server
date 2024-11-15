package tests

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"

	"github.com/vbulash/chat-server/internal/api/chat"
	"github.com/vbulash/chat-server/internal/converter"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/vbulash/chat-server/internal/model"
	"github.com/vbulash/chat-server/internal/service"
	serviceMocks "github.com/vbulash/chat-server/internal/service/mocks"
	desc "github.com/vbulash/chat-server/pkg/chat_v2"
)

func TestUpdate(t *testing.T) {
	t.Parallel()

	type chatServiceMockFunc func(mc *minimock.Controller) service.ChatService

	type args struct {
		ctx     context.Context
		request *desc.ChangeRequest
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

		request = &desc.ChangeRequest{
			Id:         id,
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
		name            string
		args            args
		err             error
		chatServiceMock chatServiceMockFunc
	}{
		{
			name: "Успешный вариант",
			args: args{
				ctx:     ctx,
				request: request,
			},
			err: nil,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.ChangeMock.Expect(ctx, id, info).Return(nil)
				return mock
			},
		},
		{
			name: "Неуспешный вариант",
			args: args{
				ctx:     ctx,
				request: request,
			},
			err: serviceErr,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.ChangeMock.Expect(ctx, id, info).Return(serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userServiceMock := tt.chatServiceMock(mc)
			api := chat.NewAPI(userServiceMock)

			_, err := api.Change(tt.args.ctx, tt.args.request)
			require.Equal(t, tt.err, err)
		})
	}
}
