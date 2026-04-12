CREATE TABLE IF NOT EXISTS urls (
    id SERIAL PRIMARY KEY,
    short_code VARCHAR(20) NOT NULL UNIQUE,
    original_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    clicks INT DEFAULT 0
    );

CREATE INDEX IF NOT EXISTS idx_short_code ON urls(short_code);