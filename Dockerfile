FROM golang:1.23.2 AS builder

COPY . /github.com/vbulash/chat_server/src/
WORKDIR /github.com/vbulash/chat_server/src/

RUN go mod download
RUN go build -o ./bin/chat_server cmd/grpc_server/main.go
# RUN go build -o ./bin/chat_client cmd/grpc_client/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/vbulash/chat_server/src/bin/chat_server .
# COPY --from=builder /github.com/vbulash/chat_server/src/bin/chat_client .
RUN ls -la .

CMD ["./chat_server"]