-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS clients (
  id SERIAL NOT NULL,
  client_id UUID NOT NULL,
  client_secret VARCHAR(255) NOT NULL,
  redirect_uri VARCHAR(255) NOT NULL,
  "app_name" VARCHAR(100) NOT NULL,
  created_at TIMESTAMPTZ NULL,
	updated_at TIMESTAMPTZ NULL,
  CONSTRAINT clients_pkey PRIMARY KEY ("id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS clients;
-- +goose StatementEnd
