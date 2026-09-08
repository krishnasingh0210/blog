CREATE EXTENSION IF NOT EXISTS pgcrypto; /*this is for random generation of UUIDs*/

CREATE TABLE users (
    id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name    VARCHAR(255) NOT NULL,
    email   VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role    VARCHAR(20) NOT NULL DEFAULT 'guest'
            CHECK (role IN ('admin', 'guest')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);