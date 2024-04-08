-- +goose Up
-- +goose StatementBegin
INSERT INTO subreddits (name, subtype, schedule) VALUES
('wallpaper', 0, '0 0 * * *'), -- every day at midnight
('wallpapers', 0, '0 0 * * *'); -- every day at midnight

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM subreddits WHERE name IN ('wallpaper', 'wallpapers');
-- +goose StatementEnd
