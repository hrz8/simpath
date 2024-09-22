-- +goose Up
-- +goose StatementBegin
ALTER TABLE access_tokens
  ADD COLUMN deleted_at TIMESTAMPTZ NULL;

ALTER TABLE refresh_tokens
  ADD COLUMN deleted_at TIMESTAMPTZ NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE access_tokens
  DROP COLUMN IF EXISTS deleted_at;

ALTER TABLE refresh_tokens
  DROP COLUMN IF EXISTS deleted_at;
-- +goose StatementEnd
