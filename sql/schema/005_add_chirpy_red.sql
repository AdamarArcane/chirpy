-- +goose Up
ALTER TABLE users
ADD is_chirpy_red BOOLEAN NOT NULL DEFAULT 'false';

-- +goose Down
ALTER TABLE users
DROP is_chirpy_red;