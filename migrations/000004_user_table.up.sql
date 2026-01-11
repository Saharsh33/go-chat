CREATE TABLE users (
    id serial primary key,

    username TEXT NOT NULL UNIQUE,
    --email TEXT UNIQUE,

    --display_name TEXT,
    --avatar_url TEXT,

    --password_hash TEXT, -- NULL if using OAuth only

    --status TEXT NOT NULL DEFAULT 'active',
    --role TEXT NOT NULL DEFAULT 'user',

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_seen_at TIMESTAMP
);