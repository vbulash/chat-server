package api

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	desc "github.com/vbulash/chat-server/pkg/chat_v2"
)

// ChatAPI Слой API
type ChatAPI interface {
	CreateSend(ctx context.Context, request *desc.CreateSendRequest) (*desc.CreateSendResponse, error)
	Get(ctx context.Context, request *desc.GetRequest) (*desc.GetResponse, error)
	Change(ctx context.Context, request *desc.ChangeRequest) (*empty.Empty, error)
	Delete(ctx context.Context, request *desc.DeleteRequest) (*empty.Empty, error)
}
