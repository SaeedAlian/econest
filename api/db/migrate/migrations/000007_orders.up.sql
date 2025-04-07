CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    total_price FLOAT NOT NULL,
    delivery_date DATE NOT NULL,
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    transaction_id INTEGER UNIQUE NOT NULL REFERENCES wallet_transactions(id) ON DELETE RESTRICT
);

CREATE TABLE order_product_variants (
    id SERIAL PRIMARY KEY,
    quantity INTEGER NOT NULL CHECK (quantity >= 1),
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    variant_id INTEGER NOT NULL REFERENCES product_variants(id) ON DELETE CASCADE
);
