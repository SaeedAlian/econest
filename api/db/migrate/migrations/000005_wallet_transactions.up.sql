CREATE TABLE wallet_transactions (
  id SERIAL PRIMARY KEY,
  amount FLOAT8 NOT NULL CHECK (amount >= 0),
  tx_type VARCHAR(20) NOT NULL,
  status VARCHAR(20) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  wallet_id INTEGER NOT NULL REFERENCES wallets(id) ON DELETE RESTRICT
);

ALTER TABLE wallet_transactions 
  ALTER COLUMN tx_type TYPE transaction_types USING tx_type::transaction_types,
  ALTER COLUMN status TYPE transaction_statuses USING status::transaction_statuses;

ALTER TABLE wallet_transactions 
  ALTER COLUMN status SET DEFAULT 'pending'::transaction_statuses;

CREATE OR REPLACE FUNCTION update_wallet_balance_on_tx_success()
RETURNS TRIGGER AS $$
DECLARE
  dl FLOAT8;
BEGIN
  IF NEW.status = 'successful' AND OLD.status = 'pending' THEN
    PERFORM 1 FROM wallets WHERE id = NEW.wallet_id FOR UPDATE;

    IF NEW.tx_type = 'deposit' THEN
      dl := NEW.amount;
    ELSIF NEW.tx_type = 'withdraw' THEN
      dl := -NEW.amount;
    ELSE
      RAISE EXCEPTION 'unknown transaction type: %', NEW.tx_type;
    END IF;

    UPDATE wallets
    SET balance = balance + dl,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.wallet_id;

  END IF;

  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION prevent_tx_update_status_after_final_state()
RETURNS TRIGGER AS $$
BEGIN
  IF OLD.status IN ('successful', 'failed') AND NEW.status IS DISTINCT FROM OLD.status THEN
    RAISE EXCEPTION 'cannot change status from % to % after it is finalized', OLD.status, NEW.status;
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_prevent_tx_status_change_after_final_state
BEFORE UPDATE ON wallet_transactions
FOR EACH ROW
WHEN (OLD.status IN ('successful', 'failed') AND NEW.status IS DISTINCT FROM OLD.status)
EXECUTE FUNCTION prevent_tx_update_status_after_final_state();

CREATE TRIGGER trg_update_wallet_balance_on_tx_success
AFTER UPDATE ON wallet_transactions
FOR EACH ROW
WHEN (OLD.status = 'pending' AND NEW.status = 'successful')
EXECUTE FUNCTION update_wallet_balance_on_tx_success();
