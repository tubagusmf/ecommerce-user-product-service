
-- +migrate Up
CREATE TABLE order_items (
    "id" SERIAL PRIMARY KEY,
    "order_id" VARCHAR(100) NOT NULL REFERENCES orders("id") ON DELETE CASCADE,
    "product_id" INT NOT NULL REFERENCES products("id") ON DELETE CASCADE,
    "quantity" INT NOT NULL,
    "price" DECIMAL NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

-- +migrate Down
DROP TABLE order_items;