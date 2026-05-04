-- Add conversation_id FK to directmessages so every DM belongs to a conversation.
ALTER TABLE directmessages
    ADD COLUMN conversation_id INT REFERENCES conversations(id) ON DELETE CASCADE;

-- Back-fill existing rows: create a conversation for each unique sender/receiver pair,
-- then link every DM to its conversation.
-- Step 1: Delete self-DMs (sender = receiver) — these are invalid.
DELETE FROM directmessages WHERE sender = receiver;

-- Step 2: Insert missing conversations (canonical ordering).
INSERT INTO conversations (user_one, user_two)
SELECT DISTINCT
    LEAST(sender, receiver),
    GREATEST(sender, receiver)
FROM directmessages
WHERE sender <> receiver
ON CONFLICT DO NOTHING;

-- Step 3: Update each DM with its conversation_id.
UPDATE directmessages d
SET conversation_id = c.id
FROM conversations c
WHERE LEAST(d.sender, d.receiver)    = c.user_one
  AND GREATEST(d.sender, d.receiver) = c.user_two;

-- Now make it NOT NULL.
ALTER TABLE directmessages ALTER COLUMN conversation_id SET NOT NULL;

-- Replace the old index with a conversation-scoped one.
DROP INDEX IF EXISTS idx_directmessages_room_created_at;

CREATE INDEX idx_directmessages_conversation_created_at
ON directmessages (conversation_id, created_at DESC);
