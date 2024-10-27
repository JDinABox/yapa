-- name: CreateUser :one
INSERT INTO Users (id, created, email, name) VALUES (?, ?, ?, ?) RETURNING *;

-- name: CleanupSessions :exec
DELETE FROM Sessions WHERE expires <= ?;