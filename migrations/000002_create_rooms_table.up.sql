CREATE TABLE rooms (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    description TEXT,
    created_by TEXT NOT NULL REFERENCES users(username),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
