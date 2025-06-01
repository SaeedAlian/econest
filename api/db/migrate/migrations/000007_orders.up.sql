CREATE TABLE orders (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT
);

CREATE TABLE order_payments (
  id SERIAL PRIMARY KEY,
  total_variants_price FLOAT8 NOT NULL CHECK (total_variants_price >= 0),
  total_shipment_price FLOAT8 NOT NULL CHECK (total_shipment_price >= 0),
  fee FLOAT8 NOT NULL CHECK (fee >= 0),
  status VARCHAR(20) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  order_id INTEGER NOT NULL UNIQUE REFERENCES orders(id) ON DELETE CASCADE
);

ALTER TABLE order_payments
  ALTER COLUMN status TYPE order_payment_statuses USING status::order_payment_statuses;

ALTER TABLE order_payments
  ALTER COLUMN status SET DEFAULT 'pending'::order_payment_statuses;

CREATE TABLE order_shipments (
  id SERIAL PRIMARY KEY,
  arrival_date DATE NOT NULL,
  status VARCHAR(30) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  order_id INTEGER NOT NULL UNIQUE REFERENCES orders(id) ON DELETE CASCADE,
  receiver_address_id INTEGER NOT NULL REFERENCES addresses(id) ON DELETE RESTRICT
);

ALTER TABLE order_shipments
  ALTER COLUMN status TYPE order_shipment_statuses USING status::order_shipment_statuses;

ALTER TABLE order_shipments
  ALTER COLUMN status SET DEFAULT 'to_be_determined'::order_shipment_statuses;

CREATE TABLE order_product_variants (
  id SERIAL PRIMARY KEY,
  quantity INTEGER NOT NULL CHECK (quantity >= 1),
  variant_price FLOAT8 NOT NULL CHECK (variant_price >= 0),
  shipping_price FLOAT8 NOT NULL CHECK (shipping_price >= 0),

  order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
  variant_id INTEGER NOT NULL REFERENCES product_variants(id) ON DELETE RESTRICT
);

CREATE OR REPLACE FUNCTION check_receiver_address_in_order_shipment()
RETURNS TRIGGER AS $$
DECLARE
  customer_id INTEGER;
BEGIN
  SELECT user_id INTO customer_id FROM orders WHERE id = NEW.order_id;

  IF NOT FOUND THEN
    RAISE EXCEPTION 'order not found for shipment: %', NEW.order_id;
  END IF;

  PERFORM 1 FROM addresses
  WHERE id = NEW.receiver_address_id AND user_id = customer_id;

  IF NOT FOUND THEN
    RAISE EXCEPTION 'receiver address % is not owned by customer %', NEW.receiver_address_id, customer_id;
  END IF;
  
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION prevent_payment_update_status_after_final_state()
RETURNS TRIGGER AS $$
BEGIN
  IF OLD.status IN ('successful', 'failed') AND NEW.status IS DISTINCT FROM OLD.status THEN
    RAISE EXCEPTION 'cannot change status from % to % after it is finalized', OLD.status, NEW.status;
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION handle_successful_order_payment()
RETURNS TRIGGER AS $$
DECLARE
  customer_wallet_id INTEGER;
  customer_wallet_balance FLOAT8;
  dl FLOAT8;

  variant_record RECORD;
  variant_current_quantity INTEGER;
  variant_store_owner_id INTEGER;
  variant_store_owner_wallet_id INTEGER;
  variant_total_price FLOAT8;
BEGIN
  IF NEW.status = 'successful' AND OLD.status = 'pending' THEN
    SELECT w.id, w.balance INTO customer_wallet_id, customer_wallet_balance
    FROM wallets w
    JOIN orders o ON o.user_id = w.user_id
    WHERE o.id = NEW.order_id
    FOR UPDATE;

    IF NOT FOUND THEN
      RAISE EXCEPTION 'customer wallet not found for order %', NEW.order_id;
    END IF;

    dl := NEW.total_variants_price + NEW.total_shipment_price + NEW.fee;

    IF customer_wallet_balance < dl THEN
      RAISE EXCEPTION 'insufficient wallet balance: required = %, available = %',
        dl, customer_wallet_balance;
    END IF;

    UPDATE wallets
    SET balance = balance - dl,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = customer_wallet_id;

    FOR variant_record IN
      SELECT opv.variant_id, opv.quantity, opv.variant_price, opv.shipping_price, pv.product_id
      FROM order_product_variants opv
      JOIN product_variants pv ON pv.id = opv.variant_id
      WHERE opv.order_id = NEW.order_id
    LOOP
      SELECT quantity INTO variant_current_quantity
      FROM product_variants
      WHERE id = variant_record.variant_id
      FOR UPDATE;

      IF variant_current_quantity < variant_record.quantity THEN
        RAISE EXCEPTION 'quantity is not enough for product: %',
          variant_record.product_id;
      END IF;

      UPDATE product_variants
      SET
        quantity = quantity - variant_record.quantity
      WHERE id = variant_record.variant_id;

      SELECT s.owner_id INTO variant_store_owner_id
      FROM store_owned_products sop
      JOIN stores s ON sop.store_id = s.id
      WHERE sop.product_id = variant_record.product_id;

      IF NOT FOUND THEN
        RAISE EXCEPTION 'store not found for product %', variant_record.product_id;
      END IF;

      SELECT id INTO variant_store_owner_wallet_id
      FROM wallets
      WHERE user_id = variant_store_owner_id
      FOR UPDATE;

      IF NOT FOUND THEN
        RAISE EXCEPTION 'wallet not found for store owner %', variant_store_owner_id;
      END IF;

      variant_total_price := variant_record.quantity * variant_record.variant_price + variant_record.shipping_price;

      UPDATE wallets
      SET balance = balance + variant_total_price,
          updated_at = CURRENT_TIMESTAMP
      WHERE user_id = variant_store_owner_id;
    END LOOP;
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION prevent_deletion_if_payment_finalized()
RETURNS TRIGGER AS $$
DECLARE
  status VARCHAR(20);
BEGIN
  SELECT status INTO status
  FROM order_payments
  WHERE order_id = (
    CASE TG_TABLE_NAME
      WHEN 'orders' THEN OLD.id
      ELSE OLD.order_id
    END
  );

  IF status IN ('successful', 'failed') THEN
    RAISE EXCEPTION 'cannot delete % with order_id % because payment has status: %',
      TG_TABLE_NAME, (CASE TG_TABLE_NAME WHEN 'orders' THEN OLD.id ELSE OLD.order_id END), status;
  END IF;

  RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_prevent_payment_status_change_after_final_state
BEFORE UPDATE ON order_payments
FOR EACH ROW
WHEN (OLD.status IN ('successful', 'failed') AND NEW.status IS DISTINCT FROM OLD.status)
EXECUTE FUNCTION prevent_payment_update_status_after_final_state();

CREATE TRIGGER trg_handle_successful_order_payment
BEFORE UPDATE ON order_payments
FOR EACH ROW
WHEN (OLD.status = 'pending' AND NEW.status = 'successful')
EXECUTE FUNCTION handle_successful_order_payment();

CREATE TRIGGER trg_prevent_order_deletion
BEFORE DELETE ON orders
FOR EACH ROW
EXECUTE FUNCTION prevent_deletion_if_payment_finalized();

CREATE TRIGGER trg_prevent_order_payment_deletion
BEFORE DELETE ON order_payments
FOR EACH ROW
EXECUTE FUNCTION prevent_deletion_if_payment_finalized();

CREATE TRIGGER trg_prevent_order_shipment_deletion
BEFORE DELETE ON order_shipments
FOR EACH ROW
EXECUTE FUNCTION prevent_deletion_if_payment_finalized();

CREATE TRIGGER trg_check_receiver_address_in_order_shipment_insertion
BEFORE INSERT ON order_shipments
FOR EACH ROW
EXECUTE FUNCTION check_receiver_address_in_order_shipment();
