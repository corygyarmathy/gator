-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id, last_fetched_at)
VALUES ( $1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetFeed :one
SELECT * FROM feeds WHERE feeds.name = $1;

-- name: GetFeedByUrl :one
SELECT * FROM feeds WHERE feeds.url = $1;

-- name: ResetFeeds :exec
DELETE FROM feeds;

-- name: GetFeeds :many
SELECT * FROM feeds;
