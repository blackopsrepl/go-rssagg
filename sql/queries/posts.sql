-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;
--

-- name: GetPostsForUser :many
SELECT posts.* FROM posts
JOIN feed_follows ON feed_follows.feed_id = posts.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at DESC
LIMIT $2;
--

-- name: GetPostsForUserByDate :many
SELECT posts.* FROM posts
JOIN feed_follows ON feed_follows.feed_id = posts.feed_id
WHERE feed_follows.user_id = $1
AND posts.published_at >= $2
ORDER BY posts.published_at DESC
LIMIT $3;
--

-- name: GetNextPostsToCrawl :many
SELECT * FROM posts
ORDER BY last_crawled_at ASC NULLS FIRST
LIMIT $1;
--

-- name: GetNextPostsToSummarize :many
SELECT * FROM posts
ORDER BY last_summarized_at ASC NULLS FIRST
LIMIT $1;
--

-- name: MarkPostCrawled :one
UPDATE posts
SET last_crawled_at = NOW()
WHERE id = $1
RETURNING *;
--

-- name: MarkPostSummarized :one
UPDATE posts
SET last_summarized_at = NOW()
WHERE id = $1
RETURNING *;
--