-- Drop trigger
DROP TRIGGER IF EXISTS update_messages_updated_at ON messages;

-- Drop indexes
DROP INDEX IF EXISTS idx_messages_reply_to;
DROP INDEX IF EXISTS idx_messages_deleted_at;
DROP INDEX IF EXISTS idx_messages_chat_created;
DROP INDEX IF EXISTS idx_messages_created_at;
DROP INDEX IF EXISTS idx_messages_user_id;
DROP INDEX IF EXISTS idx_messages_chat_id;

-- Drop table
DROP TABLE IF EXISTS messages;