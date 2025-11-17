CREATE TABLE IF NOT EXISTS metrics (
    id TEXT PRIMARY KEY,
    mtype TEXT NOT NULL CHECK (mtype IN ('сounter', 'сauge')),
    delta BIGINT,
    value DOUBLE PRECISION,
    hash TEXT
);
