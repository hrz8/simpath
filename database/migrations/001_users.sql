-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS roles (
	id BIGSERIAL NOT NULL,
	"name" VARCHAR(50) NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	CONSTRAINT roles_pkey PRIMARY KEY ("id"),
	CONSTRAINT roles_name_ukey UNIQUE ("name")
);

CREATE TABLE IF NOT EXISTS users (
	id BIGSERIAL NOT NULL,
	email VARCHAR(255) NULL UNIQUE,
	encrypted_password VARCHAR(255) NULL,
	role_id BIGSERIAL NOT NULL,
	created_at TIMESTAMPTZ NULL,
	updated_at TIMESTAMPTZ NULL,
	CONSTRAINT users_pkey PRIMARY KEY ("id"),
	CONSTRAINT users_email_ukey UNIQUE (email),
	CONSTRAINT users_role_id_fkey FOREIGN KEY (role_id) REFERENCES roles("id") ON DELETE RESTRICT ON UPDATE RESTRICT
);
CREATE INDEX idx_users_role_id ON users USING btree (role_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;
-- +goose StatementEnd
