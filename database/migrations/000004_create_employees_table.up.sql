CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS employees (
    id UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    role UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    email VARCHAR(30) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT (NOW()),
    updated_at TIMESTAMP NOT NULL DEFAULT (NOW()),
    deleted_at TIMESTAMP
);