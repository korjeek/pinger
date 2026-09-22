-- +goose Up
CREATE TABLE monitors (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    url VARCHAR(255) NOT NULL,
    tracked BOOLEAN NOT NULL DEFAULT true,
    poll_interval_sec INT NOT NULL,
    checked_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT unique_user_url UNIQUE (user_id, url)
);
CREATE INDEX idx_monitors_scheduler ON monitors(checked_at) WHERE tracked = true;

-- +goose Down
DROP TABLE IF EXISTS monitors;
