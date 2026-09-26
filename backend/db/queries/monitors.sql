-- name: CreateMonitor :one
INSERT INTO monitors (id, user_id, url, active, poll_interval_sec)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, url, active, poll_interval_sec, checked_at;

-- name: GetMonitorByUserID :many
SELECT id, user_id, url, active, poll_interval_sec, checked_at
FROM monitors
WHERE user_id = $1
ORDER BY checked_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateMonitor :exec
UPDATE monitors
SET url = $2,
    active = $3,
    poll_interval_sec = $4
WHERE id = $1;

-- name: DeleteMonitor :exec
DELETE FROM monitors
WHERE id = $1;