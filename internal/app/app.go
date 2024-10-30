package app

import (
	"context"
	"log"
	"net"

	"github.com/vbulash/chat-server/internal/closer"
	desc "github.com/vbulash/chat-server/pkg/chat_v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// App Приложение
type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

// NewApp Инициализация приложения
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Run Запуск приложения
func (app *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return app.runGRPCServer()
}

func (app *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		app.initServiceProvider,
		app.initGRPCServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *App) initServiceProvider(_ context.Context) error {
	app.serviceProvider = newServiceProvider()
	return nil
}

func (app *App) initGRPCServer(ctx context.Context) error {
	app.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	reflection.Register(app.grpcServer)

	apiLayer := app.serviceProvider.APILayer(ctx)
	desc.RegisterChatV2Server(app.grpcServer, apiLayer)

	return nil
}

func (app *App) runGRPCServer() error {
	list, err := net.Listen("tcp", app.serviceProvider.env.Address)
	if err != nil {
		return err
	}
	log.Printf("Сервер GRPC работает на %s...", app.serviceProvider.env.Address)

	err = app.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}
