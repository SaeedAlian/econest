package db_manager

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"

	"github.com/SaeedAlian/econest/api/config"
	"github.com/SaeedAlian/econest/api/types"
)

func (m *Manager) CreateOrder(p types.CreateOrderPayload) (int, error) {
	if len(p.ProductVariants) == 0 {
		return -1, types.ErrProductVariantsAreEmpty
	}

	rowId := -1
	ctx := context.Background()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return -1, err
	}

	err = tx.QueryRow("INSERT INTO orders (user_id) VALUES ($1) RETURNING id;",
		p.UserId,
	).
		Scan(&rowId)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	var totalShipmentPrice float64 = 0
	var totalVariantsPrice float64 = 0
	var orderFee float64 = 0

	variantQtyMap := make(map[int]int, len(p.ProductVariants))
	variantIds := make([]int, len(p.ProductVariants))
	for i, pv := range p.ProductVariants {
		variantQtyMap[pv.VariantId] = pv.Quantity
		variantIds[i] = pv.VariantId
	}

	variantRows, err := tx.Query(`
		SELECT
			p.id, pv.id, pv.quantity, p.shipment_factor,
			COALESCE(
				p.price * (1 - (
					SELECT discount FROM product_offers po
					WHERE po.product_id = p.id AND po.expire_at > NOW()
					LIMIT 1
				)), p.price
			) AS final_price
		FROM product_variants pv
		JOIN products p ON p.id = pv.product_id
		WHERE pv.id = ANY($1)
	`, pq.Array(variantIds))
	if err != nil {
		tx.Rollback()
		return -1, err
	}
	defer variantRows.Close()

	insertData := make([]types.OrderProductVariantInsertData, 0, len(p.ProductVariants))

	for variantRows.Next() {
		var productId int = -1
		var variantId int = -1
		var currentQuantity int = -1
		var shipmentFactor float64 = 0
		var variantPrice float64 = 0
		err := variantRows.Scan(
			&productId,
			&variantId,
			&currentQuantity,
			&shipmentFactor,
			&variantPrice,
		)
		if err != nil {
			tx.Rollback()
			return -1, err
		}
		if productId == -1 {
			tx.Rollback()
			return -1, types.ErrProductNotFound
		}

		selectedQuantity, ok := variantQtyMap[variantId]
		if !ok {
			tx.Rollback()
			return -1, types.ErrProductVariantNotFound
		}

		if currentQuantity < selectedQuantity {
			tx.Rollback()
			return -1, types.ErrProductQuantityIsNotEnough(productId)
		}

		shippingPrice := config.Env.ShipmentPrice * shipmentFactor

		totalShipmentPrice += shippingPrice
		totalVariantsPrice += variantPrice * float64(selectedQuantity)

		insertData = append(insertData, types.OrderProductVariantInsertData{
			Quantity:      selectedQuantity,
			VariantPrice:  variantPrice,
			ShippingPrice: shippingPrice,
			VariantId:     variantId,
			OrderId:       rowId,
		})
	}

	if len(insertData) != len(variantIds) {
		tx.Rollback()
		return -1, types.ErrProductVariantNotFound
	}

	for _, d := range insertData {
		_, err = tx.Exec(
			"INSERT INTO order_product_variants (quantity, variant_price, shipping_price, variant_id, order_id) VALUES ($1, $2, $3, $4, $5)",
			d.Quantity,
			d.VariantPrice,
			d.ShippingPrice,
			d.VariantId,
			d.OrderId,
		)
		if err != nil {
			tx.Rollback()
			return -1, err
		}
	}

	orderFee = totalVariantsPrice * config.Env.OrderFeeFactor

	_, err = tx.Exec(
		"INSERT INTO order_shipments (arrival_date, order_id, receiver_address_id) VALUES ($1, $2, $3)",
		p.ArrivalDate,
		rowId,
		p.ReceiverAddressId,
	)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	_, err = tx.Exec(
		"INSERT INTO order_payments (total_variants_price, total_shipment_price, fee, order_id) VALUES ($1, $2, $3, $4)",
		totalVariantsPrice,
		totalShipmentPrice,
		orderFee,
		rowId,
	)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return -1, err
	}

	return rowId, nil
}

func (m *Manager) GetOrders(
	query types.OrderSearchQuery,
) ([]types.Order, error) {
	var base string
	base = `
		SELECT
			o.*, op.status, os.status,
			op.total_variants_price, op.total_shipment_price, op.fee,
			(
				SELECT COUNT(*) 
				FROM order_product_variants opv 
				WHERE opv.order_id = o.id
		  ) AS total_products
		FROM orders o 
		JOIN order_payments op ON op.order_id = o.id
		JOIN order_shipments os ON os.order_id = o.id
	`

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

func (m *Manager) GetOrdersWithFullInfo(
	query types.OrderSearchQuery,
) ([]types.OrderWithFullInfo, error) {
	var base string
	base = `
		SELECT
			o.*, op.*, os.*, a.*,
			(
				SELECT COUNT(*) 
				FROM order_product_variants opv 
				WHERE opv.order_id = o.id
		  ) AS total_products
		FROM orders o 
		JOIN order_payments op ON op.order_id = o.id
		JOIN order_shipments os ON os.order_id = o.id
		JOIN addresses a ON a.id = os.receiver_address_id AND a.user_id IS NOT NULL
	`

	q, args := buildOrderSearchQuery(query, base)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []types.OrderWithFullInfo{}

	for rows.Next() {
		order, err := scanOrderWithFullInfoRow(rows)
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
	base = `
		SELECT COUNT(DISTINCT o.id) as count FROM orders o
		JOIN order_payments op ON op.order_id = o.id
		JOIN order_shipments os ON os.order_id = o.id
	`

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
			OrderProductVariant: *orderProductVariant,
			SelectedVariant:     *selectedVariant,
			Product:             *product,
		})
	}

	return variantsInfo, nil
}

func (m *Manager) GetOrderById(id int) (*types.Order, error) {
	rows, err := m.db.Query(`
		SELECT
			o.*, op.status, os.status,
			op.total_variants_price, op.total_shipment_price, op.fee,
			(
				SELECT COUNT(*) 
				FROM order_product_variants opv 
				WHERE opv.order_id = o.id
		  ) AS total_products
		FROM orders o 
		JOIN order_payments op ON op.order_id = o.id
		JOIN order_shipments os ON os.order_id = o.id
		WHERE o.id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	order := new(types.Order)
	order.Id = -1

	for rows.Next() {
		order, err = scanOrderRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if order.Id == -1 {
		return nil, types.ErrOrderNotFound
	}

	return order, nil
}

func (m *Manager) GetOrderWithFullInfoById(id int) (*types.OrderWithFullInfo, error) {
	rows, err := m.db.Query(`
		SELECT
			o.*, op.*, os.*, a.*,
			(
				SELECT COUNT(*) 
				FROM order_product_variants opv 
				WHERE opv.order_id = o.id
		  ) AS total_products
		FROM orders o 
		JOIN order_payments op ON op.order_id = o.id
		JOIN order_shipments os ON os.order_id = o.id
		JOIN addresses a ON a.id = os.receiver_address_id AND a.user_id IS NOT NULL
		WHERE o.id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	order := new(types.OrderWithFullInfo)
	order.Id = -1

	for rows.Next() {
		order, err = scanOrderWithFullInfoRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if order.Id == -1 {
		return nil, types.ErrOrderNotFound
	}

	return order, nil
}

func (m *Manager) IsStoreHasParticipationInOrder(orderId int, storeId int) (bool, error) {
	rows, err := m.db.Query(`
		SELECT o.id FROM orders o 
		WHERE o.id = $1 AND EXISTS (
			SELECT 1 FROM order_product_variants opv
			JOIN product_variants pv ON pv.id = opv.variant_id
			JOIN store_owned_products sop ON sop.product_id = pv.product_id
			WHERE sop.store_id = $2
		)
	`, orderId, storeId)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	id := -1

	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return false, err
		}
	}

	if id == -1 {
		return false, nil
	}

	return true, nil
}

func (m *Manager) UpdateOrderShipment(
	orderId int,
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

	if p.Status != nil {
		clauses = append(clauses, fmt.Sprintf("status = $%d", argsPos))
		args = append(args, *p.Status)
		argsPos++
	}

	if len(clauses) == 0 {
		return fmt.Errorf("No fields received to update")
	}

	args = append(args, orderId)
	q := fmt.Sprintf(
		"UPDATE order_shipments SET %s WHERE order_id = $%d",
		strings.Join(clauses, ", "),
		argsPos,
	)

	_, err := m.db.Exec(q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) UpdateOrderPayment(
	orderId int,
	p types.UpdateOrderPaymentPayload,
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
		return fmt.Errorf("No fields received to update")
	}

	args = append(args, orderId)
	q := fmt.Sprintf(
		"UPDATE order_payments SET %s WHERE order_id = $%d",
		strings.Join(clauses, ", "),
		argsPos,
	)

	_, err := m.db.Exec(q, args...)
	if err != nil {
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

func scanOrderRow(rows *sql.Rows) (*types.Order, error) {
	n := new(types.Order)

	err := rows.Scan(
		&n.Id,
		&n.CreatedAt,
		&n.UpdatedAt,
		&n.UserId,
		&n.PaymentStatus,
		&n.ShipmentStatus,
		&n.TotalVariantsPrice,
		&n.TotalShipmentPrice,
		&n.Fee,
		&n.TotalProducts,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func scanOrderWithFullInfoRow(rows *sql.Rows) (*types.OrderWithFullInfo, error) {
	n := new(types.OrderWithFullInfo)

	err := rows.Scan(
		&n.Id,
		&n.CreatedAt,
		&n.UpdatedAt,
		&n.UserId,
		&n.Payment.Id,
		&n.Payment.TotalVariantsPrice,
		&n.Payment.TotalShipmentPrice,
		&n.Payment.Fee,
		&n.Payment.Status,
		&n.Payment.CreatedAt,
		&n.Payment.UpdatedAt,
		&n.Payment.OrderId,
		&n.Shipment.Id,
		&n.Shipment.ArrivalDate,
		&n.Shipment.Status,
		&n.Shipment.CreatedAt,
		&n.Shipment.UpdatedAt,
		&n.Shipment.OrderId,
		&n.Shipment.ReceiverAddressId,
		&n.Shipment.ReceiverAddress.Id,
		&n.Shipment.ReceiverAddress.State,
		&n.Shipment.ReceiverAddress.City,
		&n.Shipment.ReceiverAddress.Street,
		&n.Shipment.ReceiverAddress.Zipcode,
		&n.Shipment.ReceiverAddress.Details,
		&n.Shipment.ReceiverAddress.IsPublic,
		&n.Shipment.ReceiverAddress.CreatedAt,
		&n.Shipment.ReceiverAddress.UpdatedAt,
		&n.Shipment.ReceiverAddress.UserId,
		new(sql.NullInt32),
		&n.TotalProducts,
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
		&n.VariantPrice,
		&n.ShippingPrice,
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
		clauses = append(clauses, fmt.Sprintf("o.user_id = $%d", argsPos))
		args = append(args, *query.UserId)
		argsPos++
	}

	if query.PaymentStatus != nil {
		clauses = append(clauses, fmt.Sprintf("op.status = $%d", argsPos))
		args = append(args, *query.PaymentStatus)
		argsPos++
	}

	if query.ShipmentStatus != nil {
		clauses = append(clauses, fmt.Sprintf("os.status = $%d", argsPos))
		args = append(args, *query.ShipmentStatus)
		argsPos++
	}

	if query.StoreId != nil {
		clauses = append(clauses, fmt.Sprintf(`
      EXISTS (
				SELECT 1 FROM order_product_variants opv
				JOIN product_variants pv ON pv.id = opv.variant_id
				JOIN store_owned_products sop ON sop.product_id = pv.product_id
				WHERE sop.store_id = $%d
			)
    `, argsPos))
		args = append(args, *query.StoreId)
		argsPos++
	}

	if query.CreatedAtLessThan != nil {
		clauses = append(clauses, fmt.Sprintf("o.created_at <= $%d", argsPos))
		args = append(args, *query.CreatedAtLessThan)
		argsPos++
	}

	if query.CreatedAtMoreThan != nil {
		clauses = append(clauses, fmt.Sprintf("o.created_at >= $%d", argsPos))
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
