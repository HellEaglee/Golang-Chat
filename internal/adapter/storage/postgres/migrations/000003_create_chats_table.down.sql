-- Drop triggers
DROP TRIGGER IF EXISTS update_chat_participants_updated_at ON chat_participants;
DROP TRIGGER IF EXISTS update_chats_updated_at ON chats;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column;

-- Drop indexes
DROP INDEX IF EXISTS idx_chat_participants_chat_user;
DROP INDEX IF EXISTS idx_chat_participants_active;
DROP INDEX IF EXISTS idx_chat_participants_user_id;
DROP INDEX IF EXISTS idx_chat_participants_chat_id;
DROP INDEX IF EXISTS idx_chats_created_at;
DROP INDEX IF EXISTS idx_chats_deleted_at;

-- Drop tables
DROP TABLE IF EXISTS chat_participants;
DROP TABLE IF EXISTS chats;

-- Drop extension if not used elsewhere
-- DROP EXTENSION IF EXISTS "uuid-ossp";