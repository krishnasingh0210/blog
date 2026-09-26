CREATE TABLE comments (
    id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    commentable_type  VARCHAR(30) NOT NULL CHECK (commentable_type IN ('blog')),
    commentable_id    UUID NOT NULL,
    author_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_comment_id UUID REFERENCES comments(id) ON DELETE CASCADE,
    content           TEXT NOT NULL,
    created_at        TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_comments_commentable ON comments(commentable_type, commentable_id);
CREATE INDEX idx_comments_parent_comment_id ON comments(parent_comment_id);