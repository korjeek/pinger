-- name: CreateSession :exec
INSERT INTO sessions (id, user_id, token_hash, expires_at)
VALUES ($1, $2, $3, $4);

-- name: GetSessionByHash :one
SELECT * FROM sessions
WHERE token_hash = $1;

-- name: RotateSession :exec
UPDATE sessions
SET revoked_at = now(),
    replaced_by = $2
WHERE id = $1;

-- name: RevokeSession :exec
UPDATE sessions
SET revoked_at = now()
WHERE id = $1;

-- name: RevokeAllUserSessions :exec
UPDATE sessions
SET revoked_at = now()
WHERE user_id = $1 AND revoked_at IS NULL;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions
WHERE expires_at < now()
   OR (revoked_at IS NOT NULL AND revoked_at < $1);
