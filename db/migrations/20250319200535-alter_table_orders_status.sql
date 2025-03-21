
-- +migrate Up
CREATE TYPE "status" AS ENUM ('pending', 'success', 'failed');

ALTER TABLE orders ADD COLUMN "status" status NOT NULL DEFAULT 'pending';

-- +migrate Down
ALTER TABLE orders DROP COLUMN "status";
