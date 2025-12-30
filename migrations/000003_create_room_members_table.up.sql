CREATE TABLE room_members (
    room_id INT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    username TEXT NOT NULL UNIQUE,
    joined_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (room_id, username)
);
