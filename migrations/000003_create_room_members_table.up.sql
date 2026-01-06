CREATE TABLE room_members (
    room_id INT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    username TEXT NOT NULL,
    room_name TEXT NOT NULL,
    joined_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (room_id, username)
);
