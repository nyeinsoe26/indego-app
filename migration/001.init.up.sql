-- Create the indego_snapshots table to store the bike availability data
CREATE TABLE indego_snapshots (
    id SERIAL PRIMARY KEY,
    timestamp TIMESTAMPTZ NOT NULL, -- UTC timestamp of the snapshot
    data JSONB NOT NULL              -- JSON data for the Indego bike stations
);

-- Create the weather_snapshots table to store the weather data
CREATE TABLE weather_snapshots (
    id SERIAL PRIMARY KEY,
    timestamp TIMESTAMPTZ NOT NULL, -- UTC timestamp of the snapshot
    data JSONB NOT NULL             -- JSON data for the weather
);

-- Create a table to link the indego and weather snapshots by timestamp
CREATE TABLE snapshots (
    id SERIAL PRIMARY KEY,
    timestamp TIMESTAMPTZ NOT NULL, -- The exact time of the snapshot
    indego_snapshot_id INT REFERENCES indego_snapshots(id) ON DELETE CASCADE,
    weather_snapshot_id INT REFERENCES weather_snapshots(id) ON DELETE CASCADE
);

-- Index on timestamp for fast queries when fetching snapshots by time
CREATE INDEX idx_indego_snapshots_timestamp ON indego_snapshots(timestamp);
CREATE INDEX idx_weather_snapshots_timestamp ON weather_snapshots(timestamp);
CREATE INDEX idx_snapshots_timestamp ON snapshots(timestamp);
