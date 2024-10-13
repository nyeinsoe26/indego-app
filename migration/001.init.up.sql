-- Create the indego_snapshots table to store the bike availability data
CREATE TABLE indego_snapshots (
    id UUID NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL, -- UTC timestamp of the snapshot
    data JSONB NOT NULL,              -- JSON data for the Indego bike stations
    CONSTRAINT indego_snapshots_pkey PRIMARY KEY (id)
);

-- Create the weather_snapshots table to store the weather data
CREATE TABLE weather_snapshots (
    id UUID NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL, -- UTC timestamp of the snapshot
    data JSONB NOT NULL,             -- JSON data for the weather
    CONSTRAINT weather_snapshots_pkey PRIMARY KEY (id)
);

-- Index on timestamp for fast queries when fetching snapshots by time
CREATE INDEX idx_indego_snapshots_timestamp ON indego_snapshots(timestamp);
CREATE INDEX idx_weather_snapshots_timestamp ON weather_snapshots(timestamp);
