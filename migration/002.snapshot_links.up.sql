-- Create a table to link the indego and weather snapshots by timestamp
CREATE TABLE snapshots (
    id UUID NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL, -- The exact time of the snapshot
    indego_snapshot_id UUID REFERENCES indego_snapshots(id) ON DELETE CASCADE,
    weather_snapshot_id UUID REFERENCES weather_snapshots(id) ON DELETE CASCADE,
    CONSTRAINT snapshots_pkey PRIMARY KEY (id)
);

-- Index on timestamp for fast queries when fetching snapshots by time
CREATE INDEX idx_snapshots_timestamp ON snapshots(timestamp);
