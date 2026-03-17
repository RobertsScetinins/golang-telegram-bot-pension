-- CREATE TABLE chats (
--     id BIGSERIAL PRIMARY KEY,
--     chat_id BIGINT NOT NULL UNIQUE,
--     type TEXT NOT NULL,
--     title TEXT,
--     created_at TIMESTAMP NOT NULL DEFAULT now()
-- );

CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL,
    message_id BIGINT NOT NULL,
    username TEXT,
    text TEXT,
    is_edited BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP,

    -- CONSTRAINT fk_chat
    --     FOREIGN KEY (chat_id)
    --     REFERENCES chats(chat_id)
    --     ON DELETE CASCADE,

    CONSTRAINT unique_message
        UNIQUE(chat_id, message_id)
);

-- CREATE INDEX idx_messages_chat_created
-- ON messages(chat_id, created_at DESC);