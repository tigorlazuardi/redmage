-- +goose Up
-- +goose StatementBegin
INSERT INTO subreddits (name, subtype, schedule) VALUES
('wallpapers', 0, '0 0 * * *'), -- every day at midnight
('wallpaper', 0, '0 0 * * *'); -- every day at midnight

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM subreddits WHERE name IN ('wallpapers', 'wallpaper');
-- +goose StatementEnd
