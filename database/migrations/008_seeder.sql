-- +goose Up
-- +goose StatementBegin
INSERT INTO roles (id, "name", created_at, updated_at) VALUES
	(1, 'root','2024-09-19 00:00:00.000','2024-09-19 00:00:00.000'),
	(2, 'user','2024-09-19 00:00:00.000','2024-09-19 00:00:00.000');
INSERT INTO scopes (id, "scope", description, is_default, created_at, updated_at) VALUES
  (1, 'read','Read permissions',true,'2024-09-19 00:00:00.000','2024-09-19 00:00:00.000'),
  (2, 'read_write','Full permissions',false,'2024-09-19 00:00:00.000','2024-09-19 00:00:00.000');
INSERT INTO clients (id, client_id,client_secret,redirect_uri,app_name,created_at,updated_at) VALUES
  (1, '600ef080-d02c-426d-bf79-64247ba0fc90','$2a$10$CUoGytf1pR7CC6Y043gt/.vFJUV4IRqvH5R6F0VfITP8s2TqrQ.4e','https://www.example.com','Test Client ABC','2024-09-19 00:00:00.000','2024-09-19 00:00:00.000');
INSERT INTO users (id, email, encrypted_password, role_id, created_at, updated_at, public_id) VALUES
	(1, 'test@root','$2a$10$4J4t9xuWhOKhfjN0bOKNReS9sL3BVSN9zxIr2.VaWWQfRBWh1dQIS',1,'2024-09-19 00:00:00.000','2024-09-19 00:00:00.000','72746765-83ea-4e90-aecd-f85a465193c9');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE id = 1;
DELETE FROM clients WHERE id = 1;
DELETE FROM scopes WHERE id = 1;
DELETE FROM scopes WHERE id = 2;
DELETE FROM roles WHERE id = 1;
DELETE FROM roles WHERE id = 2;
-- +goose StatementEnd
