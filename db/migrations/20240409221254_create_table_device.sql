-- +goose Up
-- +goose StatementBegin
CREATE TABLE devices(
    slug VARCHAR(255) NOT NULL PRIMARY KEY,
    enable INTEGER NOT NULL DEFAULT 1,
    name VARCHAR(255) NOT NULL,
    resolution_x DOUBLE NOT NULL,
    resolution_y DOUBLE NOT NULL,
    aspect_ratio_tolerance DOUBLE NOT NULL default 0.2,
    min_x INTEGER NOT NULL DEFAULT 0,
    min_y INTEGER NOT NULL DEFAULT 0,
    max_x INTEGER NOT NULL DEFAULT 0,
    max_y INTEGER NOT NULL DEFAULT 0,
    nsfw INTEGER NOT NULL DEFAULT 0,
    windows_wallpaper_mode INTEGER NOT NULL DEFAULT 0,
    created_at BIGINT DEFAULT 0 NOT NULL,
    updated_at BIGINT DEFAULT 0 NOT NULL
);

CREATE UNIQUE INDEX idx_devices_name ON devices(slug);

CREATE INDEX idx_devices_enable ON devices(enable);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS devices;
-- +goose StatementEnd
