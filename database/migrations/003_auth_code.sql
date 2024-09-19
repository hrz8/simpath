-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS authorization_codes (
	id BIGSERIAL NOT NULL,
	client_id SERIAL NOT NULL,
	"user_id" BIGSERIAL NOT NULL,
	"code" VARCHAR(50) NOT NULL,
	redirect_uri VARCHAR(255) NULL,
	"scope" VARCHAR(100) NOT NULL,
	expires_at TIMESTAMPTZ NOT NULL,
	created_at TIMESTAMPTZ NULL,
	updated_at TIMESTAMPTZ NULL,
	CONSTRAINT authorization_codes_pkey PRIMARY KEY ("id"),
	CONSTRAINT authorization_codes_code_ukey UNIQUE ("code"),
	CONSTRAINT authorization_codes_client_id_fkey FOREIGN KEY (client_id) REFERENCES clients("id") ON DELETE RESTRICT ON UPDATE RESTRICT,
	CONSTRAINT authorization_codes_user_id_fkey FOREIGN KEY ("user_id") REFERENCES users("id") ON DELETE RESTRICT ON UPDATE RESTRICT
);
CREATE INDEX idx_authorization_codes_client_id ON authorization_codes USING btree (client_id);
CREATE INDEX idx_authorization_codes_user_id ON authorization_codes USING btree ("user_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS authorization_codes;
-- +goose StatementEnd
