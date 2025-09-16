DROP TABLE IF EXISTS news CASCADE;

CREATE TABLE news (
    id SERIAL PRIMARY KEY,
    source TEXT NOT NULL,
    source_type TEXT NOT NULL,
    title TEXT NOT NULL,
    summary TEXT,
    url TEXT,
    authors TEXT[] DEFAULT '{}',
    published_at TIMESTAMPTZ NOT NULL,
    tags TEXT[] DEFAULT '{}',
    metadata JSONB DEFAULT '{}'::jsonb,
    content_hash TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for performance
CREATE UNIQUE INDEX idx_news_content_hash ON news(content_hash);
CREATE INDEX idx_news_published_at ON news(published_at DESC);
CREATE INDEX idx_news_tags ON news USING GIN (tags);
CREATE INDEX idx_news_metadata ON news USING GIN (metadata);
