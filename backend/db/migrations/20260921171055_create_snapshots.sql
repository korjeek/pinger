-- +goose Up
CREATE TABLE snapshots (
    id UUID PRIMARY KEY,
    monitor_id UUID NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
    checked_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    is_up BOOLEAN NOT NULL,
    status_code SMALLINT NOT NULL,
    response_time_ms INT NOT NULL,
    response_size BIGINT NOT NULL DEFAULT 0,
    server_name VARCHAR(100) NOT NULL DEFAULT '',
    ssl_expires_at TIMESTAMP WITH TIME ZONE
);
CREATE INDEX idx_snapshots_monitor_time ON snapshots(monitor_id, checked_at DESC);

-- +goose Down
DROP TABLE snapshots;
