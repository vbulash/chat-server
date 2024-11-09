-- +goose Up
-- +goose StatementBegin
CREATE TABLE chats (
    id SERIAL PRIMARY KEY,
    body TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE chats;
-- +goose StatementEnd