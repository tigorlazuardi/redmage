-- name: RecentlyAddedImages :many
SELECT * FROM images
WHERE created_at > ?
ORDER BY created_at
LIMIT ? OFFSET ?;