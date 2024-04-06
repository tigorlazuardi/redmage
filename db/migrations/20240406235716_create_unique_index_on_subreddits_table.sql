-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX idx_subreddits_name ON subreddits (name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_subreddits_name;
-- +goose StatementEnd
