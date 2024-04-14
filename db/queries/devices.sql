-- name: DeviceGetAll :many
SELECT * FROM devices
ORDER BY name;

-- name: DeviceCount :one
SELECT COUNT(*) FROM devices;

-- name: DeviceList :many
SELECT * FROM devices
ORDER BY name
LIMIT ? OFFSET ?;

-- name: DeviceSearch :many
SELECT * FROM devices
WHERE (name LIKE ? OR slug LIKE ?)
ORDER BY name
LIMIT ? OFFSET ?;

-- name: DeviceSearchCount :one
SELECT COUNT(*) FROM devices
WHERE (name LIKE ? OR slug LIKE ?)
ORDER BY name
LIMIT ? OFFSET ?;

-- name: DeviceCreate :one
INSERT INTO devices (name, resolution_x, resolution_y, aspect_ratio_tolerance, min_x, min_y, max_x, max_y, nsfw, windows_wallpaper_mode)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;
