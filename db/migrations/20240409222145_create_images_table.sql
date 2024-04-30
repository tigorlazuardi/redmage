-- +goose Up
-- +goose StatementBegin
CREATE TABLE images(
    id INTEGER PRIMARY KEY,
    subreddit VARCHAR(255) NOT NULL,
    device VARCHAR(250) NOT NULL,
    post_title VARCHAR(255) NOT NULL,
    post_name VARCHAR(255) NOT NULL,
    post_url VARCHAR(255) NOT NULL,
    post_created BIGINT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    post_author VARCHAR(50) NOT NULL,
    post_author_url VARCHAR(255) NOT NULL,
    image_relative_path VARCHAR(255) NOT NULL,
    image_original_url VARCHAR(255) NOT NULL,
    image_height INTEGER NOT NULL DEFAULT 0,
    image_width INTEGER NOT NULL DEFAULT 0,
    image_size BIGINT NOT NULL DEFAULT 0,
    thumbnail_relative_path VARCHAR(255) NOT NULL DEFAULT '',
    nsfw INTEGER NOT NULL DEFAULT 0,
    created_at BIGINT DEFAULT 0 NOT NULL,
    updated_at BIGINT DEFAULT 0 NOT NULL,
    CONSTRAINT fk_image_subreddit
        FOREIGN KEY (subreddit)
        REFERENCES subreddits(name)
        ON DELETE CASCADE,
    CONSTRAINT fk_image_devices_slug
        FOREIGN KEY (device)
        REFERENCES devices(slug)
        ON DELETE CASCADE
);

CREATE INDEX idx_subreddit_images ON images(subreddit);
CREATE INDEX idx_nsfw_images ON images(nsfw);
CREATE INDEX idx_images_created_at_nsfw ON images(created_at DESC, nsfw);
CREATE UNIQUE INDEX idx_unique_images_per_device ON images(device, post_name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS images;
-- +goose StatementEnd
