CREATE TABLE orders (
  id SERIAL PRIMARY KEY,
  verified BOOLEAN NOT NULL DEFAULT FALSE,
  status VARCHAR(20) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
  transaction_id INTEGER REFERENCES wallet_transactions(id) ON DELETE RESTRICT
);

ALTER TABLE orders
  ALTER COLUMN status TYPE order_statuses USING status::order_statuses;

ALTER TABLE orders 
  ALTER COLUMN status SET DEFAULT 'pending_payment'::order_statuses;

CREATE TABLE order_shipments (
  id SERIAL PRIMARY KEY,
  arrival_date DATE NOT NULL,
  shipment_date DATE NOT NULL,
  status VARCHAR(30) NOT NULL,
  shipment_type VARCHAR(30) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  order_id INTEGER NOT NULL UNIQUE REFERENCES orders(id) ON DELETE CASCADE,
  receiver_address_id INTEGER REFERENCES addresses(id) ON DELETE SET NULL,
  sender_address_id INTEGER REFERENCES addresses(id) ON DELETE SET NULL
);

ALTER TABLE order_shipments
  ALTER COLUMN shipment_type TYPE shipment_types USING shipment_type::shipment_types,
  ALTER COLUMN status TYPE shipment_statuses USING status::shipment_statuses;

CREATE TABLE order_product_variants (
  id SERIAL PRIMARY KEY,
  quantity INTEGER NOT NULL CHECK (quantity >= 1),

  order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
  variant_id INTEGER NOT NULL REFERENCES product_variants(id) ON DELETE CASCADE
);
