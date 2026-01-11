CREATE TABLE roommessages (
    id SERIAL PRIMARY KEY,
    room_id INT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    username TEXT NOT NULL REFERENCES users(username),
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_roommessages_room_id_created_at
ON roommessages (room_id, created_at DESC);
