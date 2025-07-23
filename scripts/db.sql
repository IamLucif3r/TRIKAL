CREATE TABLE articles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    title TEXT NOT NULL,
    link TEXT NOT NULL UNIQUE,
    summary_raw TEXT,
    full_text TEXT,
    source TEXT,
    
    published_at TIMESTAMP,
    fetched_at TIMESTAMP DEFAULT NOW(),

    tags TEXT[],
    usefulness_score INT CHECK (usefulness_score BETWEEN 0 AND 10),
    reason_to_post TEXT,
    is_trendworthy BOOLEAN DEFAULT FALSE
);

CREATE EXTENSION IF NOT EXISTS "pgcrypto";
