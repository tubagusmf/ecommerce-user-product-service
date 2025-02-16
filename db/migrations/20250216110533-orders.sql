
-- +migrate Up
CREATE TABLE orders (
    "id" VARCHAR(100) PRIMARY KEY,
    "user_id" INT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    "total_amount" DECIMAL NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

-- +migrate Down
DROP TABLE orders;
