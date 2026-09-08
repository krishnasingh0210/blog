CREATE EXTENSION IF NOT EXISTS pg_trgm;/*this is for full-text search*/

CREATE TABLE blogs (
    id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title   VARCHAR(255) NOT NULL,
    slug    VARCHAR(280) NOT NULL UNIQUE,
    content TEXT NOT NULL,
    status  VARCHAR(20) NOT NULL DEFAULT 'draft'
            CHECK (status IN ('draft', 'published')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_blogs_author_id ON blogs(author_id);
CREATE INDEX idx_blogs_status ON blogs(status);
CREATE INDEX idx_blogs_title ON blogs USING gin (title gin_trgm_ops);