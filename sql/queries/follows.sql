-- name: GetFollowsForUser :many
SELECT * FROM feed_follows WHERE user_id = $1;
--

-- name: CreateFollow :one
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
--

-- name: DeleteFollow :exec
-- must prevent one user from deleting follows from other users by requiring both id and user_id
DELETE FROM feed_follows WHERE id = $1 and user_id = $2;
--