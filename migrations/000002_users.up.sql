-- Users Table
CREATE TABLE users (
    id bigserial PRIMARY KEY,
    created_at timestamp WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at timestamp WITH TIME ZONE NOT NULL DEFAULT now(),
    username TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash bytea NOT NULL,
    activated BOOLEAN NOT NULL,
    VERSION INTEGER NOT NULL DEFAULT 1
);