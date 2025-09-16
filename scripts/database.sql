
-- ARTICLES ITEMS TABLE

CREATE TABLE news_items (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    url TEXT,
    source TEXT,
    published_at TIMESTAMP,
    -- LLM scoring metadata
    timeliness_score INT,
    actionability_score INT,
    explainability_score INT,
    security_relevance_score INT,
    depth_score INT,
    final_score FLOAT,
    created_at TIMESTAMP DEFAULT NOW()
);


