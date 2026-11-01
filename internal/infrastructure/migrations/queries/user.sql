-- name: CreateUser :one
INSERT INTO users(
    username,
    hash_password
)
VALUES(
    sqlc.arg(username),
    sqlc.arg(hash_password)
)
RETURNING 
    id,
    username;

-- name: Login :one
SELECT id ,username ,hash_password
FROM users
WHERE username = sqlc.arg(username);


-- name: AddPoint :one
UPDATE users
SET point = COALESCE(point, 0) + sqlc.arg(point)
WHERE id = sqlc.arg(id)
RETURNING username, point;

-- name: AddESPoint :exec
INSERT INTO events(
    user_id,
    amount
)
VALUES (
    sqlc.arg(user_id)::integer,
    sqlc.arg(amount)::integer
);

-- name: LeaderBoard :many
SELECT username,point
FROM users
ORDER BY point DESC;
