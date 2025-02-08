
-- +migrate Up
CREATE TYPE "roles" AS ENUM ('admin', 'customer', 'seller');

CREATE TABLE users (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "email" VARCHAR(255) NOT NULL UNIQUE,
    "password" VARCHAR(255) NOT NULL,
    "role" roles NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

-- +migrate Down

DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS roles;