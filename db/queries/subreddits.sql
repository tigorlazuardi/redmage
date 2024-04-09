-- name: SubredditsGetAll :many
SELECT * FROM subreddits;

-- name: SubredditsList :many
SELECT * FROM subreddits
ORDER BY name
LIMIT ? OFFSET ?;

-- name: SubredditsListCount :one
SELECT COUNT(*) From subreddits;

-- name: SubredditsSearch :many
SELECT * FROM subreddits
WHERE name LIKE ?
ORDER BY name
LIMIT ? OFFSET ?;

-- name: SubredditsSearchCount :one
SELECT COUNT(*) FROM subreddits
WHERE name LIKE ?
ORDER BY name
LIMIT ? OFFSET ?;

-- name: SubredditCreate :one
INSERT INTO subreddits (name, subtype, schedule)
VALUES (?, ?, ?)
RETURNING *;

