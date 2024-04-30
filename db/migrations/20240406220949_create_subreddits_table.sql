-- +goose Up
-- +goose StatementBegin
CREATE TABLE subreddits (
    name VARCHAR(30) NOT NULL PRIMARY KEY,
    enable_schedule INT NOT NULL DEFAULT 1,
    subtype INT NOT NULL DEFAULT 0,
    schedule VARCHAR(20) NOT NULL DEFAULT '@daily',
    countback INT NOT NULL DEFAULT 100,
    created_at BIGINT DEFAULT 0 NOT NULL,
    updated_at BIGINT DEFAULT 0 NOT NULL
);

CREATE UNIQUE INDEX idx_subreddits_name ON subreddits (name);

CREATE TRIGGER update_subreddits_timestamp AFTER UPDATE ON subreddits FOR EACH ROW
BEGIN
    UPDATE subreddits SET updated_at = CURRENT_TIMESTAMP WHERE id = old.id;
END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE subreddits;
-- +goose StatementEnd
