-- Restore old index
CREATE INDEX IF NOT EXISTS idx_directmessages_room_created_at
ON directmessages (sender, receiver, created_at DESC);

DROP INDEX IF EXISTS idx_directmessages_conversation_created_at;

ALTER TABLE directmessages DROP COLUMN IF EXISTS conversation_id;
