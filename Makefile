include .env

LOCAL_BIN:=$(CURDIR)/bin
MIGRATION_DIR=./database/migrations
MIGRATION_DSN="host=$(DB_HOST) port=$(DB_PORT) dbname=$(DB_DATABASE) user=$(DB_USERNAME) password=$(DB_PASSWORD) sslmode=disable"

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.22.1

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml --fix

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate:
	make generate-chat-api

generate-chat-api:
	mkdir -p pkg/chat_v2
	protoc --proto_path grpc/api/chat/v2 \
	--go_out=pkg/chat_v2 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/chat_v2 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	--experimental_allow_proto3_optional \
	grpc/api/chat/v2/chat.proto

migration-status:
	$(LOCAL_BIN)/goose -dir $(MIGRATION_DIR) postgres $(MIGRATION_DSN) status -v

migration-migrate:
	$(LOCAL_BIN)/goose -dir $(MIGRATION_DIR) postgres $(MIGRATION_DSN) up -v

migration-rollback:
	$(LOCAL_BIN)/goose -dir $(MIGRATION_DIR) postgres $(MIGRATION_DSN) down -v

test:
	go clean -testcache
	go test ./... -covermode count -coverpkg=github.com/vbulash/chat-server/internal/service/...,github.com/vbulash/chat-server/internal/api/... -count 5

test-coverage:
	go clean -testcache
	go test ./... -coverprofile=coverage.tmp.out -covermode count -coverpkg=github.com/vbulash/chat-server/internal/service/...,github.com/vbulash/chat-server/internal/api/... -count 5
	grep -v 'mocks\|config' coverage.tmp.out  > coverage.out
	rm coverage.tmp.out
	go tool cover -html=coverage.out;
	go tool cover -func=./coverage.out | grep "total";
	grep -sqFx "/coverage.out" .gitignore || echo "/coverage.out" >> .gitignore
