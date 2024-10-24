FROM golang:1.23.2 AS builder

COPY . /github.com/vbulash/chat_server/src/
WORKDIR /github.com/vbulash/chat_server/src/

RUN go mod download
RUN go build -o ./bin/chat_server cmd/grpc_server/main.go
# RUN go build -o ./bin/chat_client cmd/grpc_client/main.go

FROM alpine:3.20.3

WORKDIR /root/
COPY --from=builder /github.com/vbulash/chat_server/src/bin/chat_server .
# COPY --from=builder /github.com/vbulash/chat_server/src/bin/chat_client .
RUN chmod +x ./chat_server

CMD ["./chat_server"]