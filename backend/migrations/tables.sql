CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users(
    guid UUID PRIMARY KEY, 
    username TEXT,
    password TEXT,
    deauthorized BOOLEAN DEFAULT FALSE
);
CREATE TABLE refresh_tokens(
    token TEXT PRIMARY KEY,
    user_guid UUID REFERENCES users(guid),
    issued_at TIMESTAMP,
    expires_at TIMESTAMP,
    valid BOOLEAN DEFAULT TRUE,
    user_agent TEXT
) 