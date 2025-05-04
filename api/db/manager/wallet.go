package db_manager

import (
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

	var q string

	if len(clauses) == 0 {
		q = "SELECT * FROM wallet_transactions"
	} else {
		q = fmt.Sprintf("SELECT * FROM wallet_transactions WHERE %s", strings.Join(clauses, " AND "))
	}

	if query.Offset != nil {
		q = fmt.Sprintf("%s OFFSET $%d", q, argsPos)
		args = append(args, *query.Offset)
		argsPos++
	}

	if query.Limit != nil {
		q = fmt.Sprintf("%s LIMIT $%d", q, argsPos)
		args = append(args, *query.Limit)
		argsPos++
	}

	q = fmt.Sprintf("%s;", q)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, err
	}

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

func (m *Manager) UpdateWallet(id int, p types.UpdateWalletPayload) error {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if p.Balance != nil {
		clauses = append(clauses, fmt.Sprintf("balance = $%d", argsPos))
		args = append(args, *p.Balance)
		argsPos++
	}

	if len(clauses) == 0 {
		return types.ErrNoFieldsReceivedToUpdate
	}

	clauses = append(clauses, fmt.Sprintf("updated_at = $%d", argsPos))
	args = append(args, time.Now())
	argsPos++

	args = append(args, id)
	q := fmt.Sprintf(
		"UPDATE wallets SET %s WHERE id = $%d",
		strings.Join(clauses, ", "),
		argsPos,
	)

	_, err := m.db.Exec(q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) UpdateWalletTransaction(id int, p types.UpdateWalletTransactionPayload) error {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if p.Status != nil {
		clauses = append(clauses, fmt.Sprintf("status = $%d", argsPos))
		args = append(args, *p.Status)
		argsPos++
	}

	if len(clauses) == 0 {
		return types.ErrNoFieldsReceivedToUpdate
	}

	clauses = append(clauses, fmt.Sprintf("updated_at = $%d", argsPos))
	args = append(args, time.Now())
	argsPos++

	args = append(args, id)
	q := fmt.Sprintf(
		"UPDATE product_tags SET %s WHERE id = $%d",
		strings.Join(clauses, ", "),
		argsPos,
	)

	_, err := m.db.Exec(q, args...)
	if err != nil {
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
