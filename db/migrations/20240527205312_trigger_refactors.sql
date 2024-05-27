-- +goose Up
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_subreddits_timestamp; -- Faulty trigger. Must be removed and never recovered.

CREATE TRIGGER subreddits_update_timestamp_on_update AFTER UPDATE ON subreddits FOR EACH ROW
BEGIN
    UPDATE subreddits SET updated_at = unixepoch() WHERE name = old.name;
END;

CREATE TRIGGER devices_update_timestamp_on_update AFTER UPDATE ON devices FOR EACH ROW
BEGIN
    UPDATE devices SET updated_at = unixepoch() WHERE slug = old.slug;
END;

CREATE TRIGGER subreddits_update_timestamp_on_image_insert AFTER INSERT ON images FOR EACH ROW
BEGIN
    UPDATE subreddits SET updated_at = unixepoch() WHERE name = new.subreddit; -- new -> image row.
END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS subreddits_update_timestamp_on_update;
DROP TRIGGER IF EXISTS devices_update_timestamp_on_update;
DROP TRIGGER IF EXISTS subreddits_update_timestamp_on_image_insert;
-- +goose StatementEnd
