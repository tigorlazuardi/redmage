-- +goose Up
-- +goose StatementBegin
CREATE TABLE devices(
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    resolution_x DOUBLE NOT NULL,
    resolution_y DOUBLE NOT NULL,
    aspect_ratio_tolerance DOUBLE NOT NULL default 0.2,
    min_x INTEGER NOT NULL DEFAULT 0,
    min_y INTEGER NOT NULL DEFAULT 0,
    max_x INTEGER NOT NULL DEFAULT 0,
    max_y INTEGER NOT NULL DEFAULT 0,
    nsfw INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_devices_name ON devices(name);

CREATE TRIGGER update_devices_timestamp AFTER UPDATE ON devices FOR EACH ROW
BEGIN
    UPDATE devices SET updated_at = CURRENT_TIMESTAMP WHERE id = old.id;
END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS devices;
-- +goose StatementEnd
