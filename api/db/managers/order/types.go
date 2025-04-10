package order_db_manager

import "time"

type OrderStatus string

const (
	OrderStatusPendingPayment  OrderStatus = "pending_payment"
	OrderStatusPendingDelivery OrderStatus = "pending_delivery"
	OrderStatusDelivered       OrderStatus = "delivered"
	OrderStatusCancelled       OrderStatus = "cancelled"
)

var ValidOrderStatuses = []OrderStatus{
	OrderStatusPendingPayment,
	OrderStatusPendingDelivery,
	OrderStatusDelivered,
	OrderStatusCancelled,
}

func (s OrderStatus) IsValid() bool {
	for _, v := range ValidOrderStatuses {
		if s == v {
			return true
		}
	}

	return false
}

type Order struct {
	Id            int
	TotalPrice    float32
	DeliveryDate  time.Time
	Verified      bool
	Status        OrderStatus
	CreatedAt     time.Time
	UpdatedAt     time.Time
	UserId        int
	TransactionId int
}

type OrderProductVariant struct {
	Id        int
	Quantity  int
	OrderId   int
	VariantId int
}
