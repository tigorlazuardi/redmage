-- +goose Up
-- +goose StatementBegin
CREATE TABLE blacklists(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    device VARCHAR(255) NOT NULL COLLATE NOCASE DEFAULT '',
    subreddit VARCHAR(255) NOT NULL COLLATE NOCASE,
    post_name VARCHAR(255) NOT NULL,
    created_at BIGINT DEFAULT 0 NOT NULL,
    CONSTRAINT fk_blacklist_subreddit_post_name
        FOREIGN KEY (subreddit, post_name)
        REFERENCES images(subreddit, post_name)
        ON DELETE CASCADE
);

CREATE INDEX idx_blacklist_device ON blacklists(device);
CREATE INDEX idx_blacklist_subreddit_post_name ON blacklists(subreddit, post_name);
CREATE INDEX idx_blacklist_created_at ON blacklists(created_at DESC);
CREATE UNIQUE INDEX idx_unique_blacklist ON blacklists(device, subreddit, post_name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS blacklists;
-- +goose StatementEnd
