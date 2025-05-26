CREATE TABLE stores (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) UNIQUE NOT NULL,
  description VARCHAR(1023) NOT NULL,
  -- FIX: must turn the default to false
  -- TODO: create a verification system for stores
  verified BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  owner_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE stores_settings (
  id SERIAL PRIMARY KEY,
  public_owner BOOLEAN NOT NULL DEFAULT FALSE,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  store_id INTEGER NOT NULL REFERENCES stores(id) ON DELETE CASCADE
);

ALTER TABLE addresses
  ADD COLUMN store_id INTEGER REFERENCES stores(id) ON DELETE CASCADE;

ALTER TABLE addresses
  ALTER COLUMN user_id DROP NOT NULL;

ALTER TABLE addresses
  ADD CONSTRAINT address_has_owner CHECK (
    (user_id IS NOT NULL AND store_id IS NULL) OR
    (user_id IS NULL AND store_id IS NOT NULL)
  );

ALTER TABLE phonenumbers
  ADD COLUMN store_id INTEGER REFERENCES stores(id) ON DELETE CASCADE;

ALTER TABLE phonenumbers
  ALTER COLUMN user_id DROP NOT NULL;

ALTER TABLE phonenumbers
  ADD CONSTRAINT phonenumber_has_owner CHECK (
    (user_id IS NOT NULL AND store_id IS NULL) OR
    (user_id IS NULL AND store_id IS NOT NULL)
  );

CREATE TABLE store_owned_products (
  store_id INTEGER NOT NULL REFERENCES stores(id) ON DELETE RESTRICT,
  product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  PRIMARY KEY (store_id, product_id)
);
