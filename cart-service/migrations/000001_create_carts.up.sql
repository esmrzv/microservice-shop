CREATE EXTENSION IF NOT EXISTS "pgcrypto";


CREATE TABLE carts(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE TABLE cart_items(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cart_id UUID NOT NULL,
    product_id UUID NOT NULL,
    quantity INT NOT NULL CHECK ( quantity > 0 ),

    CONSTRAINT fk_cart_items_cart
                       FOREIGN KEY (cart_id)
                       REFERENCES carts(id)
                       ON DELETE CASCADE,
    CONSTRAINT uq_cart_product
                       UNIQUE (cart_id, product_id)
);