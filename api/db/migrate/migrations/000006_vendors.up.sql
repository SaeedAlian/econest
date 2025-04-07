CREATE TABLE vendors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    description VARCHAR(1023) NOT NULL,
    verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    owner_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

ALTER TABLE addresses
ADD COLUMN vendor_id INTEGER REFERENCES vendors(id) ON DELETE CASCADE;

ALTER TABLE addresses
ALTER COLUMN user_id DROP NOT NULL;

ALTER TABLE addresses
ADD CONSTRAINT address_has_owner CHECK (
    (user_id IS NOT NULL AND vendor_id IS NULL) OR
    (user_id IS NULL AND vendor_id IS NOT NULL)
);

ALTER TABLE phonenumbers
ADD COLUMN vendor_id INTEGER REFERENCES vendors(id) ON DELETE CASCADE;

ALTER TABLE phonenumbers
ALTER COLUMN user_id DROP NOT NULL;

ALTER TABLE phonenumbers
ADD CONSTRAINT phonenumber_has_owner CHECK (
    (user_id IS NOT NULL AND vendor_id IS NULL) OR
    (user_id IS NULL AND vendor_id IS NOT NULL)
);

CREATE TABLE vendor_products (
    vendor_id INTEGER NOT NULL REFERENCES vendors(id) ON DELETE RESTRICT,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    PRIMARY KEY (vendor_id, product_id)
);
