CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE orders(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL ,
    status VARCHAR(50) NOT NULL,
    total_price NUMERIC(12, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE order_items(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL,
    product_id UUID NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    price NUMERIC(12, 2) NOT NULL CHECK ( price >= 0 ),

    CONSTRAINT fk_order_items_order
                        FOREIGN KEY (order_id)
                        REFERENCES orders(id)
                        ON DELETE CASCADE
);