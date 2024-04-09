-- +goose Up
-- +goose StatementBegin
CREATE TABLE subreddits (
    id INTEGER PRIMARY KEY,
    name VARCHAR(30) NOT NULL,
    subtype INT NOT NULL,
    schedule VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
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
