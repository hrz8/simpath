-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS access_tokens (
	id BIGSERIAL NOT NULL,
	client_id SERIAL NOT NULL,
	"user_id" BIGSERIAL NOT NULL,
	access_token VARCHAR(50) NOT NULL,
	"scope" VARCHAR(100) NOT NULL,
	expires_at TIMESTAMPTZ NOT NULL,
	created_at TIMESTAMPTZ NULL,
	updated_at TIMESTAMPTZ NULL,
	CONSTRAINT access_tokens_pkey PRIMARY KEY ("id"),
	CONSTRAINT access_tokens_token_ukey UNIQUE (access_token),
	CONSTRAINT access_tokens_client_id_fkey FOREIGN KEY (client_id) REFERENCES clients("id") ON DELETE RESTRICT ON UPDATE RESTRICT,
	CONSTRAINT access_tokens_user_id_fkey FOREIGN KEY ("user_id") REFERENCES users("id") ON DELETE RESTRICT ON UPDATE RESTRICT
);
CREATE INDEX idx_access_tokens_client_id ON access_tokens USING btree (client_id);
CREATE INDEX idx_access_tokens_user_id ON access_tokens USING btree ("user_id");

CREATE TABLE IF NOT EXISTS refresh_tokens (
	id BIGSERIAL NOT NULL,
	client_id SERIAL NOT NULL,
	"user_id" BIGSERIAL NOT NULL,
	refresh_token VARCHAR(40) NOT NULL,
	"scope" VARCHAR(200) NOT NULL,
	expires_at timestamptz NOT NULL,
	created_at TIMESTAMPTZ NULL,
	updated_at TIMESTAMPTZ NULL,
	CONSTRAINT refresh_tokens_pkey PRIMARY KEY ("id"),
	CONSTRAINT refresh_tokens_token_ukey UNIQUE (refresh_token),
	CONSTRAINT refresh_tokens_client_id_fkey FOREIGN KEY (client_id) REFERENCES clients("id") ON DELETE RESTRICT ON UPDATE RESTRICT,
	CONSTRAINT refresh_tokens_user_id_fkey FOREIGN KEY ("user_id") REFERENCES users("id") ON DELETE RESTRICT ON UPDATE RESTRICT
);
CREATE INDEX idx_refresh_tokens_client_id ON refresh_tokens USING btree (client_id);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens USING btree ("user_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS access_tokens;
-- +goose StatementEnd
