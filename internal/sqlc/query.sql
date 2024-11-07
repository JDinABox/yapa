-- name: CreateUser :one
INSERT INTO Users (id, created, email, name) VALUES (?, ?, ?, ?) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM Users WHERE id = ? LIMIT 1;

-- name: CreateSession :one
INSERT INTO Sessions (id, created, updated, expires, data) VALUES (?, ?, ?, ?, ?) RETURNING *;

-- name: GetSession :one
Select * from Sessions where id = ? LIMIT 1;

-- name: UpdateSessionData :exec
UPDATE Sessions SET updated = ?, expires = ?, data = ? WHERE id = ?;

-- name: DeleteSession :exec
DELETE FROM Sessions WHERE id = ?;

-- name: CleanupSessions :exec
DELETE FROM Sessions WHERE expires <= ?;