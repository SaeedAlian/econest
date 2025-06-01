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
		"INSERT INTO wallet_transactions (amount, tx_type, wallet_id) VALUES ($1, $2, $3) RETURNING id;",
		p.Amount,
		p.TxType,
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

func (m *Manager) GetWalletTransactionById(id int) (*types.WalletTransaction, error) {
	rows, err := m.db.Query(
		"SELECT * FROM wallet_transactions WHERE id = $1;",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tx := new(types.WalletTransaction)
	tx.Id = -1

	for rows.Next() {
		tx, err = scanWalletTransactionRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if tx.Id == -1 {
		return nil, types.ErrWalletTransactionNotFound
	}

	return tx, nil
}

func (m *Manager) GetUserWallet(userId int) (*types.Wallet, error) {
	rows, err := m.db.Query(
		"SELECT * FROM wallets WHERE user_id = $1;",
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	wallet := new(types.Wallet)
	wallet.Id = -1

	for rows.Next() {
		wallet, err = scanWalletRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if wallet.Id == -1 {
		return nil, types.ErrWalletNotFound
	}

	return wallet, nil
}

func (m *Manager) UpdateWallet(id int, p types.UpdateWalletPayload) error {
	clauses := []string{}
	args := []any{}
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

func (m *Manager) UpdateWalletTransaction(
	walletId int,
	transactionId int,
	p types.UpdateWalletTransactionPayload,
) error {
	clauses := []string{}
	args := []any{}
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

	args = append(args, walletId)
	args = append(args, transactionId)
	q := fmt.Sprintf(
		"UPDATE wallet_transactions SET %s WHERE wallet_id = $%d AND id = $%d",
		strings.Join(clauses, ", "),
		argsPos,
		argsPos+1,
	)

	_, err := m.db.Exec(q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) DeleteWalletTransaction(walletId int, transactionId int) error {
	_, err := m.db.Exec(
		"DELETE FROM wallet_transactions WHERE id = $1 AND wallet_id = $2;",
		transactionId, walletId,
	)
	if err != nil {
		return err
	}

	return nil
}

func scanWalletRow(rows *sql.Rows) (*types.Wallet, error) {
	n := new(types.Wallet)

	err := rows.Scan(
		&n.Id,
		&n.Balance,
		&n.CreatedAt,
		&n.UpdatedAt,
		&n.UserId,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func scanWalletTransactionRow(rows *sql.Rows) (*types.WalletTransaction, error) {
	n := new(types.WalletTransaction)

	err := rows.Scan(
		&n.Id,
		&n.Amount,
		&n.TxType,
		&n.Status,
		&n.CreatedAt,
		&n.UpdatedAt,
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
) (string, []any) {
	clauses := []string{}
	args := []any{}
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
		clauses = append(
			clauses,
			fmt.Sprintf("wallet_id IN (SELECT id FROM wallets WHERE user_id = $%d)", argsPos),
		)
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
