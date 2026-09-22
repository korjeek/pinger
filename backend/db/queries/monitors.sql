-- name: CreateWebsite :one
INSERT INTO monitors (id, user_id, url, tracked, poll_interval_sec)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, url, tracked, poll_interval_sec, checked_at;

-- name: GetWebsitesByUserID :many
SELECT id, user_id, url, tracked, poll_interval_sec, checked_at
FROM monitors
WHERE user_id = $1;

-- name: UpdateWebsite :exec
UPDATE monitors
SET url = $2,
    tracked = $3,
    poll_interval_sec = $4
WHERE id = $1;

-- name: DeleteWebsite :exec
DELETE FROM monitors
WHERE id = $1;