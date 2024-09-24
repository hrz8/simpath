-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
  ADD COLUMN public_id UUID NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
  DROP COLUMN IF EXISTS public_id;
-- +goose StatementEnd
