
-- ARTICLES ITEMS TABLE

CREATE TABLE IF NOT EXISTS articles (
    id SERIAL PRIMARY KEY,
    title TEXT,
    description TEXT,
    link TEXT UNIQUE,
    published_at TIMESTAMP,
    source TEXT
);
ALTER TABLE articles ADD COLUMN llm_score REAL DEFAULT 0;

