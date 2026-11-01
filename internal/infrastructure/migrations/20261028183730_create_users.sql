-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    hash_password TEXT NOT NULL,
    point INTEGER DEFAULT 0
);


-- +goose Down
DROP TABLE users;
