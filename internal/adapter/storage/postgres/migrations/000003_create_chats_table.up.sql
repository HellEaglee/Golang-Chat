-- Create chats table
CREATE TABLE IF NOT EXISTS chats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100), -- Optional chat name (for group chats)
    is_group BOOLEAN NOT NULL DEFAULT false, -- true for group chats, false for direct messages
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Create chat_participants table (many-to-many relationship)
CREATE TABLE IF NOT EXISTS chat_participants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    chat_id UUID NOT NULL,
    user_id UUID NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'member', -- 'admin', 'member', 'moderator'
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    left_at TIMESTAMPTZ, -- When user left the chat
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    
    -- Explicit foreign key constraints
    CONSTRAINT fk_chat_participants_chat_id FOREIGN KEY (chat_id) REFERENCES chats(id),
    CONSTRAINT fk_chat_participants_user_id FOREIGN KEY (user_id) REFERENCES users(id),
    UNIQUE(chat_id, user_id) -- Prevent duplicate participants
);

-- Indexes for chats
CREATE INDEX IF NOT EXISTS idx_chats_deleted_at ON chats(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_chats_created_at ON chats(created_at);

-- Indexes for chat_participants
CREATE INDEX IF NOT EXISTS idx_chat_participants_chat_id ON chat_participants(chat_id);
CREATE INDEX IF NOT EXISTS idx_chat_participants_user_id ON chat_participants(user_id);
CREATE INDEX IF NOT EXISTS idx_chat_participants_active ON chat_participants(is_active) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_chat_participants_chat_user ON chat_participants(chat_id, user_id);
CREATE INDEX IF NOT EXISTS idx_chat_participants_deleted_at ON chat_participants(deleted_at) WHERE deleted_at IS NULL;

-- Trigger for updating updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers for automatic updated_at
CREATE TRIGGER update_chats_updated_at 
    BEFORE UPDATE ON chats 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_chat_participants_updated_at 
    BEFORE UPDATE ON chat_participants 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();