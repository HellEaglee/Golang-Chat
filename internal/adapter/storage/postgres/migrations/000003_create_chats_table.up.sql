CREATE TABLE IF NOT EXISTS chats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100), -- Optional chat name (for group chats)
    is_group BOOLEAN NOT NULL DEFAULT false, -- true for group chats, false for direct messages
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS chat_participants (
    chat_id UUID NOT NULL,
    user_id UUID NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'member', -- 'admin', 'member', 'moderator'
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    left_at TIMESTAMPTZ, -- When user left the chat
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    PRIMARY KEY (chat_id, user_id),
    
    CONSTRAINT fk_chat_participants_chat_id FOREIGN KEY (chat_id) REFERENCES chats(id),
    CONSTRAINT fk_chat_participants_user_id FOREIGN KEY (user_id) REFERENCES users(id)
);