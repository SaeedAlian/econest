DROP TRIGGER IF EXISTS trg_prevent_payment_status_change_after_final_state ON order_payments;
DROP FUNCTION IF EXISTS prevent_payment_update_status_after_final_state;

DROP TRIGGER IF EXISTS trg_handle_successful_order_payment ON order_payments;
DROP FUNCTION IF EXISTS handle_successful_order_payment;

DROP TRIGGER IF EXISTS trg_prevent_order_deletion ON orders;
DROP TRIGGER IF EXISTS trg_prevent_order_payment_deletion ON order_payments;
DROP TRIGGER IF EXISTS trg_prevent_order_shipment_deletion ON order_shipments;
DROP FUNCTION IF EXISTS prevent_deletion_if_payment_finalized;

DROP TRIGGER IF EXISTS trg_check_receiver_address_in_order_shipment_insertion ON order_shipments;
DROP FUNCTION IF EXISTS check_receiver_address_in_order_shipment;

DROP TABLE order_product_variants;
DROP TABLE order_shipments;
DROP TABLE order_payments;
DROP TABLE orders;
