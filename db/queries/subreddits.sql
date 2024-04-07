-- name: ListSubreddits :many
SELECT * FROM subreddits
ORDER BY name
LIMIT ?;

-- name: CreateSubreddit :one
INSERT INTO subreddits (name, subtype, schedule)
VALUES (?, ?, ?)
RETURNING *;

