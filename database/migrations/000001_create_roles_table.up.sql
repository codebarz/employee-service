CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
    title TEXT NOT NULL,
    level INT NOT NULL DEFAULT (0),
    created_at TIMESTAMP NOT NULL DEFAULT (NOW()),
    updated_at TIMESTAMP NOT NULL DEFAULT (NOW()),
    deleted_at TIMESTAMP
)
