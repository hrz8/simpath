-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS scopes (
	id BIGSERIAL NOT NULL,
	"scope" VARCHAR(100) NOT NULL,
	"description" TEXT NULL,
	is_default bool NULL DEFAULT false,
	created_at TIMESTAMPTZ NULL,
	updated_at TIMESTAMPTZ NULL,
	CONSTRAINT scopes_pkey PRIMARY KEY ("id"),
	CONSTRAINT scopes_scope_ukey UNIQUE ("scope")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS scopes;
-- +goose StatementEnd
