-- name: CreateSnapshot :one
INSERT INTO snapshots (
    id, monitor_id, alive, status_code, response_time_ms, response_size, server_name, ssl_expires_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetSnapshotsByMonitorID :many
SELECT id, monitor_id, checked_at, alive, status_code, response_time_ms, response_size, server_name, ssl_expires_at
FROM snapshots
WHERE monitor_id = $1
ORDER BY checked_at DESC
LIMIT $2 OFFSET $3;