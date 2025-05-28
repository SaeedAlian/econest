package db_manager

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/SaeedAlian/econest/api/types"
)

func (m *Manager) CreateOrder(p types.CreateOrderPayload) (int, error) {
	rowId := -1
	ctx := context.Background()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return -1, err
	}
	err = tx.QueryRow("INSERT INTO orders (user_id, transaction_id, status) VALUES ($1, $2, $3) RETURNING id;",
		p.UserId, p.TransactionId, types.OrderStatusPendingPayment,
	).
		Scan(&rowId)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	for _, pv := range p.ProductVariants {
		_, err = tx.Exec(
			"INSERT INTO order_product_variants (quantity, variant_id, order_id) VALUES ($1, $2, $3);",
			pv.Quantity,
			pv.VariantId,
			rowId,
		)
		if err != nil {
			tx.Rollback()
			return -1, err
		}
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return -1, err
	}

	return rowId, nil
}

func (m *Manager) CreateOrderShipment(p types.CreateOrderShipmentPayload) (int, error) {
	rowId := -1
	err := m.db.QueryRow(
		`INSERT INTO order_shipments (
      arrival_date,
      shipment_date,
      shipment_type,
      order_id,
      receiver_address_id,
      sender_address_id,
      status
    ) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`,
		p.ArrivalDate,
		p.ShipmentDate,
		p.ShipmentType,
		p.OrderId,
		p.ReceiverAddressId,
		p.SenderAddressId,
		types.ShipmentStatusToBeDetermined,
	).
		Scan(&rowId)
	if err != nil {
		return -1, err
	}

	return rowId, nil
}

func (m *Manager) GetOrders(
	query types.OrderSearchQuery,
) ([]types.Order, error) {
	var base string
	base = "SELECT * FROM orders"

	q, args := buildOrderSearchQuery(query, base)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []types.Order{}

	for rows.Next() {
		order, err := scanOrderRow(rows)
		if err != nil {
			return nil, err
		}

		orders = append(orders, *order)
	}

	return orders, nil
}

func (m *Manager) GetOrdersCount(
	query types.OrderSearchQuery,
) (int, error) {
	var base string
	base = "SELECT COUNT(*) as count FROM orders"

	q, args := buildOrderSearchQuery(query, base)

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

func (m *Manager) GetOrderShipments(orderId int) ([]types.OrderShipment, error) {
	rows, err := m.db.Query("SELECT * FROM order_shipments WHERE order_id = $1;", orderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	shipments := []types.OrderShipment{}

	for rows.Next() {
		ship, err := scanOrderShipmentRow(rows)
		if err != nil {
			return nil, err
		}

		shipments = append(shipments, *ship)
	}

	return shipments, nil
}

func (m *Manager) GetOrderProductVariants(orderId int) ([]types.OrderProductVariant, error) {
	rows, err := m.db.Query("SELECT * FROM order_product_variants WHERE order_id = $1;", orderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	variants := []types.OrderProductVariant{}

	for rows.Next() {
		v, err := scanOrderProductVariantRow(rows)
		if err != nil {
			return nil, err
		}

		variants = append(variants, *v)
	}

	return variants, nil
}

func (m *Manager) GetOrderProductVariantsInfo(
	orderId int,
) ([]types.OrderProductVariantInfo, error) {
	rows, err := m.db.Query("SELECT * FROM order_product_variants WHERE order_id = $1;", orderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	variantsInfo := []types.OrderProductVariantInfo{}

	for rows.Next() {
		orderProductVariant, err := scanOrderProductVariantRow(rows)
		if err != nil {
			return nil, err
		}

		selectedVariant, err := m.GetProductVariantWithAttributeSetById(
			orderProductVariant.VariantId,
		)
		if err != nil {
			return nil, err
		}

		product, err := m.GetProductById(selectedVariant.ProductId)
		if err != nil {
			return nil, err
		}

		variantsInfo = append(variantsInfo, types.OrderProductVariantInfo{
			Id:              orderProductVariant.Id,
			Quantity:        orderProductVariant.Quantity,
			SelectedVariant: *selectedVariant,
			Product:         *product,
		})
	}

	return variantsInfo, nil
}

func (m *Manager) UpdateOrder(
	id int,
	p types.UpdateOrderPayload,
) error {
	q, args, err := buildOrderUpdateQuery(id, p)
	if err != nil {
		return err
	}

	_, err = m.db.Exec(q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) UpdateOrderShipment(
	id int,
	p types.UpdateOrderShipmentPayload,
) error {
	clauses := []string{}
	args := []any{}
	argsPos := 1

	if p.ArrivalDate != nil {
		clauses = append(clauses, fmt.Sprintf("arrival_date = $%d", argsPos))
		args = append(args, *p.ArrivalDate)
		argsPos++
	}

	if p.ReceiverAddressId != nil {
		clauses = append(clauses, fmt.Sprintf("receiver_address_id = $%d", argsPos))
		args = append(args, *p.ReceiverAddressId)
		argsPos++
	}

	if p.SenderAddressId != nil {
		clauses = append(clauses, fmt.Sprintf("sender_address_id = $%d", argsPos))
		args = append(args, *p.SenderAddressId)
		argsPos++
	}

	if p.Status != nil {
		clauses = append(clauses, fmt.Sprintf("status = $%d", argsPos))
		args = append(args, *p.Status)
		argsPos++
	}

	if len(clauses) == 0 {
		return fmt.Errorf("No fields received to update")
	}

	args = append(args, id)
	q := fmt.Sprintf(
		"UPDATE order_shipments SET %s WHERE id = $%d",
		strings.Join(clauses, ", "),
		argsPos,
	)

	_, err := m.db.Exec(q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) UpdateOrderAndTransactionAndWallet(
	orderId int,
	orderPayload types.UpdateOrderPayload,
	walletPayload types.UpdateWalletPayload,
	transactionPayload types.UpdateWalletTransactionPayload,
) error {
	ctx := context.Background()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	transactionId := -1
	err = tx.QueryRow(
		"SELECT transaction_id FROM orders WHERE id = $1;",
		orderId,
	).Scan(&transactionId)
	if err != nil {
		tx.Rollback()
		return err
	}

	if transactionId == -1 {
		return types.ErrWalletTransactionNotFound
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

	err = updateOrderAsDBTx(tx, orderId, orderPayload)
	if err != nil {
		tx.Rollback()
		return err
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

func (m *Manager) DeleteOrder(id int) error {
	_, err := m.db.Exec(
		"DELETE FROM orders WHERE id = $1;",
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) DeleteOrderShipment(id int) error {
	_, err := m.db.Exec(
		"DELETE FROM order_shipments WHERE id = $1;",
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func scanOrderRow(rows *sql.Rows) (*types.Order, error) {
	n := new(types.Order)

	err := rows.Scan(
		&n.Id,
		&n.Verified,
		&n.Status,
		&n.CreatedAt,
		&n.UpdatedAt,
		&n.UserId,
		&n.TransactionId,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func scanOrderShipmentRow(rows *sql.Rows) (*types.OrderShipment, error) {
	n := new(types.OrderShipment)

	err := rows.Scan(
		&n.Id,
		&n.ArrivalDate,
		&n.ShipmentDate,
		&n.Status,
		&n.ShipmentType,
		&n.CreatedAt,
		&n.UpdatedAt,
		&n.OrderId,
		&n.ReceiverAddressId,
		&n.SenderAddressId,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func scanOrderProductVariantRow(rows *sql.Rows) (*types.OrderProductVariant, error) {
	n := new(types.OrderProductVariant)

	err := rows.Scan(
		&n.Id,
		&n.Quantity,
		&n.OrderId,
		&n.VariantId,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func buildOrderSearchQuery(
	query types.OrderSearchQuery,
	base string,
) (string, []any) {
	clauses := []string{}
	args := []any{}
	argsPos := 1

	if query.UserId != nil {
		clauses = append(clauses, fmt.Sprintf("user_id = $%d", argsPos))
		args = append(args, *query.UserId)
		argsPos++
	}

	if query.Verified != nil {
		clauses = append(clauses, fmt.Sprintf("verified = $%d", argsPos))
		args = append(args, *query.Verified)
		argsPos++
	}

	if query.Status != nil {
		clauses = append(clauses, fmt.Sprintf("status = $%d", argsPos))
		args = append(args, *query.Status)
		argsPos++
	}

	if query.CreatedAtLessThan != nil {
		clauses = append(clauses, fmt.Sprintf("created_at <= $%d", argsPos))
		args = append(args, *query.CreatedAtLessThan)
		argsPos++
	}

	if query.CreatedAtMoreThan != nil {
		clauses = append(clauses, fmt.Sprintf("created_at >= $%d", argsPos))
		args = append(args, *query.CreatedAtMoreThan)
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

func updateOrderAsDBTx(
	tx *sql.Tx,
	id int,
	p types.UpdateOrderPayload,
) error {
	q, args, err := buildOrderUpdateQuery(id, p)
	if err != nil {
		return err
	}

	_, err = tx.Exec(q, args...)
	if err != nil {
		return err
	}

	return nil
}

func buildOrderUpdateQuery(
	orderId int,
	p types.UpdateOrderPayload,
) (string, []any, error) {
	clauses := []string{}
	args := []any{}
	argsPos := 1

	if p.Verified != nil {
		clauses = append(clauses, fmt.Sprintf("verified = $%d", argsPos))
		args = append(args, *p.Verified)
		argsPos++
	}

	if p.Status != nil {
		clauses = append(clauses, fmt.Sprintf("status = $%d", argsPos))
		args = append(args, *p.Status)
		argsPos++
	}

	if len(clauses) == 0 {
		return "", nil, types.ErrNoFieldsReceivedToUpdate
	}

	args = append(args, orderId)
	q := fmt.Sprintf(
		"UPDATE orders SET %s WHERE id = $%d",
		strings.Join(clauses, ", "),
		argsPos,
	)

	return q, args, nil
}
