CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(32) UNIQUE NOT NULL,
    first_name VARCHAR(255) DEFAULT NULL,
    last_name VARCHAR(255) DEFAULT NULL,
    birth_date DATE CHECK(AGE(birth_date)::INT BETWEEN 0 AND 150) DEFAULT NULL,
);

CREATE INDEX IF NOT EXISTS username_idx ON users(username);

CREATE TABLE IF NOT EXISTS chat_members (
    chat_id BEGIN NOT NULL,
    user_id BEGIN NOT NULL REFERENCES users(id)
);
