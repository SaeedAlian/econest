DROP TRIGGER IF EXISTS trg_update_wallet_balance_on_tx_success ON wallet_transactions;
DROP FUNCTION IF EXISTS update_wallet_balance_on_tx_success;
DROP TRIGGER IF EXISTS trg_prevent_tx_status_change_after_final_state ON wallet_transactions;
DROP FUNCTION IF EXISTS prevent_tx_update_status_after_final_state;
DROP TRIGGER IF EXISTS trg_prevent_tx_deletion ON wallet_transactions;
DROP FUNCTION IF EXISTS prevent_deletion_if_tx_finalized;
DROP TABLE wallet_transactions;
