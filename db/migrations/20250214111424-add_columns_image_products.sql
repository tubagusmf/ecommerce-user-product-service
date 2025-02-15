
-- +migrate Up
ALTER TABLE products ADD COLUMN image_url VARCHAR(255);

-- +migrate Down

ALTER TABLE products DROP COLUMN image_url;
