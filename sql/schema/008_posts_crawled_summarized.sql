-- +goose Up
ALTER TABLE posts ADD COLUMN last_crawled_at TIMESTAMP;
ALTER TABLE posts ADD COLUMN last_summarized_at TIMESTAMP;

-- +goose Down
ALTER TABLE posts DROP COLUMN last_summarized_at;
ALTER TABLE posts DROP COLUMN last_crawled_at;
