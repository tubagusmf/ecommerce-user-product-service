
-- +migrate Up
CREATE TABLE user_sessions (
    "id" SERIAL PRIMARY KEY,
    "user_id" INT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    "token" TEXT NOT NULL,
    "expires_at" TIMESTAMP NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down

DROP TABLE IF EXISTS user_sessions;
