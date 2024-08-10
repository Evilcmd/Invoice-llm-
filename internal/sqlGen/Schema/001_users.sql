-- +goose Up
CREATE TABLE users(id UUID PRIMARY KEY, user_name TEXT UNIQUE NOT NULL, passwd TEXT NOT NULL);

-- +goose Down
DROP TABLE users;