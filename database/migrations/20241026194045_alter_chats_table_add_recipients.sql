-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.chats ADD recipients json NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.chats DROP COLUMN recipients;
-- +goose StatementEnd
