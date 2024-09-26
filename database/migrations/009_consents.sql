-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS consents (
	id BIGSERIAL NOT NULL,
	client_id SERIAL NOT NULL,
	"user_id" BIGSERIAL NOT NULL,
	consent bool NULL DEFAULT false,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
    deleted_at timestamptz NULL,
	CONSTRAINT consents_pkey PRIMARY KEY ("id"),
	CONSTRAINT consents_client_id_fkey FOREIGN KEY (client_id) REFERENCES clients("id") ON DELETE RESTRICT ON UPDATE RESTRICT,
	CONSTRAINT consents_user_id_fkey FOREIGN KEY ("user_id") REFERENCES users("id") ON DELETE RESTRICT ON UPDATE RESTRICT
);
CREATE INDEX idx_consents_client_id ON consents USING btree (client_id);
CREATE INDEX idx_consents_user_id ON consents USING btree ("user_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_consents_client_id;
DROP INDEX IF EXISTS idx_consents_user_id;
DROP TABLE IF EXISTS consents;
-- +goose StatementEnd
