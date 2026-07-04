CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(32) UNIQUE NOT NULL,
    email VARCHAR(64) UNIQUE DEFAULT NULL,
    first_name VARCHAR(100) DEFAULT NULL,
    last_name VARCHAR(100) DEFAULT NULL,
    birth_date DATE DEFAULT NULL,
    gender VARCHAR(5) CHECK(gender IN ('man', 'woman')) DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS chats (
    id BIGSERIAL PRIMARY KEY,
    type VARCHAR(16) NOT NULL
    name VARCHAR(64) NOT NULL,
    title VARCHAR(64) DEFAULT NULL,
    description VARCHAR(255)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    owner_id BIGINT DEFAULT NULL REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS chats_owner_id_idx ON chats(owner_id);

CREATE TABLE IF NOT EXISTS chat_members (
    chat_id BIGINT NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (chat_id, user_id)
);

CREATE INDEX IF NOT EXISTS chat_members_user_id_idx ON chat_members(user_id);
