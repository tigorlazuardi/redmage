-- +goose Up
-- +goose StatementBegin
CREATE TABLE images(
    id INTEGER PRIMARY KEY,
    subreddit_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    post_id VARCHAR(50) NOT NULL,
    poster VARCHAR(50) NOT NULL,
    image_relative_path VARCHAR(255) NOT NULL,
    thumbnail_relative_path VARCHAR(255) NOT NULL DEFAULT '',
    image_original_url VARCHAR(255) NOT NULL,
    thumbnail_original_url VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_subreddit_id
        FOREIGN KEY (subreddit_id)
        REFERENCES subreddits(id)
        ON DELETE CASCADE
);

CREATE TRIGGER update_images_timestamp AFTER UPDATE ON images FOR EACH ROW
BEGIN
    UPDATE images SET updated_at = CURRENT_TIMESTAMP WHERE id = old.id;
END;

CREATE INDEX idx_subreddit_id ON images(subreddit_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS images;
-- +goose StatementEnd
