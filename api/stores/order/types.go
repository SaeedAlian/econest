package order_store

import "time"

type Order struct {
	Id            int
	TotalPrice    float32
	DeliveryDate  time.Time
	Verified      bool
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
