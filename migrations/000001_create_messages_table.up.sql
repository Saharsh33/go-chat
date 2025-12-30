CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    room TEXT NOT NULL,
    username TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_messages_room_created_at
ON messages (room, created_at DESC);
