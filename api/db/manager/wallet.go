package db_manager

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/SaeedAlian/econest/api/types"
)

func (m *Manager) CreateWalletTransaction(p types.CreateWalletTransactionPayload) (int, error) {
	rowId := -1
	err := m.db.QueryRow(
		"INSERT INTO wallet_transactions (amount, tx_type, status, wallet_id) VALUES ($1, $2, $3, $4) RETURNING id;",
		p.Amount,
		p.TxType,
		types.TransactionStatusPending,
		p.WalletId,
	).
		Scan(&rowId)
	if err != nil {
		return -1, err
	}

	return rowId, nil
}

func (m *Manager) GetWalletTransactions(
	query types.WalletTransactionSearchQuery,
) ([]types.WalletTransaction, error) {
	var base string
	base = "SELECT * FROM wallet_transactions"

	q, args := buildWalletTransactionSearchQuery(query, base)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	txs := []types.WalletTransaction{}

	for rows.Next() {
		tx, err := scanWalletTransactionRow(rows)
		if err != nil {
			return nil, err
		}

		txs = append(txs, *tx)
	}

	return txs, nil
}

func (m *Manager) GetWalletTransactionsCount(
	query types.WalletTransactionSearchQuery,
) (int, error) {
	var base string
	base = "SELECT COUNT(*) as count FROM wallet_transactions"

	q, args := buildWalletTransactionSearchQuery(query, base)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return -1, err
		}
	}

	return count, nil
}

func (m *Manager) UpdateWallet(id int, p types.UpdateWalletPayload) error {
	q, args, err := buildWalletUpdateQuery(id, p)
	if err != nil {
		return err
	}

	_, err = m.db.Exec(q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) UpdateWalletTransaction(id int, p types.UpdateWalletTransactionPayload) error {
	q, args, err := buildWalletTransactionUpdateQuery(id, p)
	if err != nil {
		return err
	}

	_, err = m.db.Exec(q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) UpdateWalletAndTransaction(
	transactionId int,
	walletPayload types.UpdateWalletPayload,
	transactionPayload types.UpdateWalletTransactionPayload,
) error {
	ctx := context.Background()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	walletId := -1
	err = tx.QueryRow(
		"SELECT wallet_id FROM wallet_transactions WHERE id = $1;",
		transactionId,
	).Scan(&walletId)
	if err != nil {
		tx.Rollback()
		return err
	}

	if walletId == -1 {
		return types.ErrWalletNotFound
	}

	err = updateWalletAsDBTx(tx, walletId, walletPayload)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = updateWalletTransactionAsDBTx(tx, transactionId, transactionPayload)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func scanWalletTransactionRow(rows *sql.Rows) (*types.WalletTransaction, error) {
	n := new(types.WalletTransaction)

	err := rows.Scan(
		&n.Id,
		&n.Amount,
		&n.TxType,
		&n.Status,
		&n.CreatedAt,
		&n.WalletId,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func buildWalletTransactionSearchQuery(
	query types.WalletTransactionSearchQuery,
	base string,
) (string, []interface{}) {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if query.TxType != nil {
		clauses = append(clauses, fmt.Sprintf("tx_type = $%d", argsPos))
		args = append(args, *query.TxType)
		argsPos++
	}

	if query.Status != nil {
		clauses = append(clauses, fmt.Sprintf("status = $%d", argsPos))
		args = append(args, *query.Status)
		argsPos++
	}

	if query.BeforeDate != nil {
		clauses = append(clauses, fmt.Sprintf("created_at <= $%d", argsPos))
		args = append(args, *query.BeforeDate)
		argsPos++
	}

	if query.AfterDate != nil {
		clauses = append(clauses, fmt.Sprintf("created_at >= $%d", argsPos))
		args = append(args, *query.AfterDate)
		argsPos++
	}

	if query.UserId != nil {
		clauses = append(clauses, fmt.Sprintf("user_id = $%d", argsPos))
		args = append(args, *query.UserId)
		argsPos++
	}

	q := base
	if len(clauses) > 0 {
		q += " WHERE " + strings.Join(clauses, " AND ")
	}

	if query.Offset != nil {
		q += fmt.Sprintf(" OFFSET $%d", argsPos)
		args = append(args, *query.Offset)
		argsPos++
	}

	if query.Limit != nil {
		q += fmt.Sprintf(" LIMIT $%d", argsPos)
		args = append(args, *query.Limit)
		argsPos++
	}

	q += ";"
	return q, args
}

func updateWalletAsDBTx(tx *sql.Tx, id int, p types.UpdateWalletPayload) error {
	q, args, err := buildWalletUpdateQuery(id, p)
	if err != nil {
		return err
	}

	_, err = tx.Exec(q, args...)
	if err != nil {
		return err
	}

	return nil
}

func updateWalletTransactionAsDBTx(
	tx *sql.Tx,
	id int,
	p types.UpdateWalletTransactionPayload,
) error {
	q, args, err := buildWalletTransactionUpdateQuery(id, p)
	if err != nil {
		return err
	}

	_, err = tx.Exec(q, args...)
	if err != nil {
		return err
	}

	return nil
}

func buildWalletUpdateQuery(
	walletId int,
	p types.UpdateWalletPayload,
) (string, []interface{}, error) {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if p.Balance != nil {
		clauses = append(clauses, fmt.Sprintf("balance = $%d", argsPos))
		args = append(args, *p.Balance)
		argsPos++
	}

	if len(clauses) == 0 {
		return "", nil, types.ErrNoFieldsReceivedToUpdate
	}

	clauses = append(clauses, fmt.Sprintf("updated_at = $%d", argsPos))
	args = append(args, time.Now())
	argsPos++

	args = append(args, walletId)
	q := fmt.Sprintf(
		"UPDATE wallets SET %s WHERE id = $%d",
		strings.Join(clauses, ", "),
		argsPos,
	)

	return q, args, nil
}

func buildWalletTransactionUpdateQuery(
	transactionId int,
	p types.UpdateWalletTransactionPayload,
) (string, []interface{}, error) {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if p.Status != nil {
		clauses = append(clauses, fmt.Sprintf("status = $%d", argsPos))
		args = append(args, *p.Status)
		argsPos++
	}

	if len(clauses) == 0 {
		return "", nil, types.ErrNoFieldsReceivedToUpdate
	}

	clauses = append(clauses, fmt.Sprintf("updated_at = $%d", argsPos))
	args = append(args, time.Now())
	argsPos++

	args = append(args, transactionId)
	q := fmt.Sprintf(
		"UPDATE wallet_transactions SET %s WHERE id = $%d",
		strings.Join(clauses, ", "),
		argsPos,
	)

	return q, args, nil
}
