-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    hash_password TEXT NOT NULL,
    point INTEGER
);


-- +goose Down
DROP TABLE users;
