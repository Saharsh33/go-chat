CREATE TABLE directmessages (
    id SERIAL PRIMARY KEY,
    sender TEXT NOT NULL REFERENCES users(username),
    receiver TEXT NOT NULL REFERENCES users(username),
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_directmessages_room_created_at
ON directmessages (sender,receiver, created_at DESC);