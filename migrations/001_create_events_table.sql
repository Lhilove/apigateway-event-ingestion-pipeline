CREATE TABLE IF NOT EXISTS events (
    id          BIGSERIAL PRIMARY KEY,
    source      VARCHAR(255) NOT NULL,
    host        VARCHAR(255) NOT NULL,
    event_type  VARCHAR(100) NOT NULL,
    severity    VARCHAR(20)  NOT NULL,
    timestamp   TIMESTAMPTZ  NOT NULL,
    message     TEXT,
    metadata    JSONB,
    received_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_events_host ON events(host);
CREATE INDEX idx_events_event_type ON events(event_type);
CREATE INDEX idx_events_severity ON events(severity);
CREATE INDEX idx_events_timestamp ON events(timestamp);