-- +goose Up
-- +goose StatementBegin
CREATE TABLE schedule_status(
    id INTEGER PRIMARY KEY,
    subreddit VARCHAR(255) NOT NULL,
    status TINYINT NOT NULL DEFAULT 0,
    error_message VARCHAR(255) NOT NULL DEFAULT '',
    created_at BIGINT DEFAULT 0 NOT NULL,
    updated_at BIGINT DEFAULT 0 NOT NULL,
    CONSTRAINT fk_scheduler_status_subreddit
        FOREIGN KEY (subreddit)
        REFERENCES subreddits(name)
        ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_unique_schedule_status_per_subreddit ON schedule_status(subreddit);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS schedule_status;
-- +goose StatementEnd
