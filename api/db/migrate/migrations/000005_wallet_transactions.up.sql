CREATE TABLE wallet_transactions (
  id SERIAL PRIMARY KEY,
  amount FLOAT8 NOT NULL CHECK (amount >= 0),
  tx_type VARCHAR(20) NOT NULL,
  status VARCHAR(20) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  wallet_id INTEGER NOT NULL REFERENCES wallets(id) ON DELETE CASCADE
);

ALTER TABLE wallet_transactions 
  ALTER COLUMN tx_type TYPE transaction_types USING tx_type::transaction_types,
  ALTER COLUMN status TYPE transaction_statuses USING status::transaction_statuses;

ALTER TABLE wallet_transactions 
  ALTER COLUMN status SET DEFAULT 'pending'::transaction_statuses;
