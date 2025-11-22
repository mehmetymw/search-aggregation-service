CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TABLE IF NOT EXISTS provider_format_metadata (
    id VARCHAR(50) PRIMARY KEY,          -- json, xml, csv, rss, yaml...
    display_name VARCHAR(100) NOT NULL,
    is_enabled BOOLEAN NOT NULL DEFAULT true,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS providers (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(100) NOT NULL UNIQUE,
    format VARCHAR(50) NOT NULL REFERENCES provider_format_metadata(id),
    base_url TEXT NOT NULL,
    is_enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS content_type_metadata (
    id VARCHAR(50) PRIMARY KEY,
    display_name VARCHAR(100) NOT NULL,
    is_enabled BOOLEAN NOT NULL DEFAULT true,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS contents (
    id BIGSERIAL PRIMARY KEY,
    provider_id BIGINT NOT NULL REFERENCES providers(id),
    provider_content_id VARCHAR(255) NOT NULL,
    title TEXT NOT NULL,
    content_type VARCHAR(50) NOT NULL REFERENCES content_type_metadata(id),
    published_at TIMESTAMP NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(provider_id, provider_content_id)
);

CREATE TABLE IF NOT EXISTS content_stats (
    id BIGSERIAL PRIMARY KEY,
    content_id BIGINT NOT NULL REFERENCES contents(id) ON DELETE CASCADE,
    views BIGINT NOT NULL DEFAULT 0,
    likes BIGINT NOT NULL DEFAULT 0,
    duration_sec INTEGER NOT NULL DEFAULT 0,
    reading_time INTEGER NOT NULL DEFAULT 0,
    reactions BIGINT NOT NULL DEFAULT 0,
    comments BIGINT NOT NULL DEFAULT 0,
    last_sync_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(content_id)
);

CREATE TABLE IF NOT EXISTS tags (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS content_tags (
    content_id BIGINT NOT NULL REFERENCES contents(id) ON DELETE CASCADE,
    tag_id BIGINT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (content_id, tag_id)
);

CREATE TABLE IF NOT EXISTS content_raw_payloads (
    content_id BIGINT PRIMARY KEY REFERENCES contents(id) ON DELETE CASCADE,
    provider_id BIGINT NOT NULL REFERENCES providers(id),
    raw_payload JSONB NOT NULL,
    fetched_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS provider_sync_runs (
    id BIGSERIAL PRIMARY KEY,
    provider_id BIGINT NOT NULL REFERENCES providers(id),
    started_at TIMESTAMP NOT NULL,
    finished_at TIMESTAMP,
    status VARCHAR(50) NOT NULL,
    item_count INTEGER NOT NULL DEFAULT 0,
    error_message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_contents_type ON contents (content_type);
CREATE INDEX IF NOT EXISTS idx_contents_published ON contents (published_at DESC);
CREATE INDEX IF NOT EXISTS idx_contents_title_trgm ON contents USING gin (title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_content_stats_views ON content_stats (views DESC);
CREATE INDEX IF NOT EXISTS idx_tags_name ON tags (name);
CREATE INDEX IF NOT EXISTS idx_sync_runs_provider ON provider_sync_runs (provider_id, created_at DESC);

-- Create scoring_rules table
CREATE TABLE IF NOT EXISTS scoring_rules (
    key VARCHAR(255) PRIMARY KEY,
    value JSONB NOT NULL,
    description TEXT,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- Seed Providers (Idempotent: ON CONFLICT DO NOTHING)

INSERT INTO provider_format_metadata (id, display_name, is_enabled, sort_order) VALUES
('json', 'JSON Format', true, 1),
('xml', 'XML Format', true, 2)
ON CONFLICT (id) DO NOTHING;

INSERT INTO providers (name, code, format, base_url, is_enabled) VALUES
('json-provider', 'json-provider', 'json', 'https://raw.githubusercontent.com/WEG-Technology/mock/refs/heads/main/v2/provider1', true),
('xml-provider', 'xml-provider', 'xml', 'https://raw.githubusercontent.com/WEG-Technology/mock/refs/heads/main/v2/provider2', true)
ON CONFLICT (code) DO NOTHING;

-- Seed Scoring Rules
-- Video Configuration
INSERT INTO scoring_rules (key, value, description) VALUES
('video_config', '{
    "type_multiplier": 1.5,
    "engagement_weight": 10.0,
    "views_divisor": 1000.0,
    "likes_divisor": 100.0
}', 'Configuration for Video content scoring'),

-- Article Configuration
('article_config', '{
    "type_multiplier": 1.0,
    "engagement_weight": 5.0,
    "reading_time_divisor": 1.0,
    "reactions_divisor": 50.0
}', 'Configuration for Article content scoring'),

-- Recency Configuration
('recency_config', '{
    "week_score": 5.0,
    "month_score": 3.0,
    "quarter_score": 1.0
}', 'Configuration for Recency scoring')
ON CONFLICT (key) DO NOTHING;

-- Seed Content Type Metadata
INSERT INTO content_type_metadata (id, display_name, is_enabled, sort_order) VALUES
('video', 'Video', true, 1),
('article', 'Article', true, 2)
ON CONFLICT (id) DO NOTHING;
