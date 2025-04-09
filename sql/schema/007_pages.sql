-- +goose Up
CREATE TABLE pages (
    id UUID PRIMARY KEY,
    content TEXT,
    -- if post is deleted, page is also deleted
    -- if page is deleted, post is not affected
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE
);

-- +goose down
DROP TABLE pages;