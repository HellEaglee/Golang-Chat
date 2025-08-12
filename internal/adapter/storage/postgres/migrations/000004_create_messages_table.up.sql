-- Create messages table
CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    chat_id UUID NOT NULL,
    user_id UUID NOT NULL,
    text TEXT NOT NULL,
    is_edited BOOLEAN NOT NULL DEFAULT false,
    reply_to_message_id UUID, -- For message replies
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    
    -- Foreign key constraints
    CONSTRAINT fk_messages_chat_id FOREIGN KEY (chat_id) REFERENCES chats(id),
    CONSTRAINT fk_messages_user_id FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_messages_reply_to FOREIGN KEY (reply_to_message_id) REFERENCES messages(id)
);

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_messages_chat_id ON messages(chat_id);
CREATE INDEX IF NOT EXISTS idx_messages_user_id ON messages(user_id);
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);
CREATE INDEX IF NOT EXISTS idx_messages_chat_created ON messages(chat_id, created_at);
CREATE INDEX IF NOT EXISTS idx_messages_deleted_at ON messages(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_messages_reply_to ON messages(reply_to_message_id) WHERE reply_to_message_id IS NOT NULL;

-- Trigger for updating updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger for automatic updated_at
CREATE TRIGGER update_messages_updated_at 
    BEFORE UPDATE ON messages 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();