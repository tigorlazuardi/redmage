-- +goose Up
-- +goose StatementBegin
CREATE TABLE schedule_histories(
    id INTEGER PRIMARY KEY,
    subreddit VARCHAR(255) NOT NULL,
    status TINYINT NOT NULL DEFAULT 0,
    error_message VARCHAR(255) NOT NULL DEFAULT '',
    created_at BIGINT DEFAULT 0 NOT NULL,
    CONSTRAINT fk_scheduler_histories_subreddit
        FOREIGN KEY (subreddit)
        REFERENCES subreddits(name)
        ON DELETE CASCADE
);

CREATE INDEX idx_schedule_histories_subreddit_created_at ON schedule_histories(subreddit, created_at DESC);
CREATE INDEX idx_schedule_histories_created_at ON schedule_histories(created_at DESC);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS schedule_histories;
-- +goose StatementEnd
