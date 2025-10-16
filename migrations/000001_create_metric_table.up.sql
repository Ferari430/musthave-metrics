CREATE TABLE IF NOT EXISTS metrics (
    id TEXT PRIMARY KEY,
    mtype TEXT NOT NULL CHECK (mtype IN ('Counter', 'Gauge')),
    delta BIGINT,
    value DOUBLE PRECISION,
    hash TEXT
);
