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
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE subreddits;
-- +goose StatementEnd
