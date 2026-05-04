-- conversations table: groups DMs between two users into a virtual "thread"
-- The canonical pair is stored with user_one < user_two (lexicographic) to avoid duplicates.
CREATE TABLE conversations (
    id SERIAL PRIMARY KEY,
    user_one TEXT NOT NULL REFERENCES users(username),
    user_two TEXT NOT NULL REFERENCES users(username),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (user_one, user_two),
    CHECK (user_one < user_two)
);

CREATE INDEX idx_conversations_user_one ON conversations (user_one);
CREATE INDEX idx_conversations_user_two ON conversations (user_two);
