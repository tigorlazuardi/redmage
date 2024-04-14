-- name: GetDevices :many
SELECT * FROM devices;

-- name: CreateDevice :one
INSERT INTO devices (name, resolution_x, resolution_y, aspect_ratio_tolerance, min_x, min_y, max_x, max_y, nsfw, windows_wallpaper_mode)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;
