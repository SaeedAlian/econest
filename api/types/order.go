package types

import "time"

type OrderBase struct {
	Id        int       `json:"id"        exposure:"private,needPermission"`
	CreatedAt time.Time `json:"createdAt" exposure:"private,needPermission"`
	UpdatedAt time.Time `json:"updatedAt" exposure:"private,needPermission"`
	UserId    int       `json:"userId"    exposure:"private,needPermission"`
}

type OrderPayment struct {
	Id                 int                `json:"id"                 exposure:"private,needPermission"`
	TotalVariantsPrice float64            `json:"totalVariantsPrice" exposure:"private,needPermission"`
	TotalShipmentPrice float64            `json:"totalShipmentPrice" exposure:"private,needPermission"`
	Fee                float64            `json:"fee"                exposure:"private,needPermission"`
	Status             OrderPaymentStatus `json:"status"             exposure:"private,needPermission"`
	CreatedAt          time.Time          `json:"createdAt"          exposure:"private,needPermission"`
	UpdatedAt          time.Time          `json:"updatedAt"          exposure:"private,needPermission"`
	OrderId            int                `json:"orderId"            exposure:"private,needPermission"`
}

type OrderShipment struct {
	Id                int                 `json:"id"                exposure:"private,needPermission"`
	ArrivalDate       time.Time           `json:"arrivalDate"       exposure:"private,needPermission"`
	Status            OrderShipmentStatus `json:"status"            exposure:"private,needPermission"`
	CreatedAt         time.Time           `json:"createdAt"         exposure:"private,needPermission"`
	UpdatedAt         time.Time           `json:"updatedAt"         exposure:"private,needPermission"`
	OrderId           int                 `json:"orderId"           exposure:"private,needPermission"`
	ReceiverAddressId int                 `json:"receiverAddressId" exposure:"private,needPermission"`
}

type Order struct {
	OrderBase
	PaymentStatus      OrderPaymentStatus  `json:"paymentStatus"      exposure:"private,needPermission"`
	ShipmentStatus     OrderShipmentStatus `json:"shipmentStatus"     exposure:"private,needPermission"`
	TotalVariantsPrice float64             `json:"totalVariantsPrice" exposure:"private,needPermission"`
	TotalShipmentPrice float64             `json:"totalShipmentPrice" exposure:"private,needPermission"`
	Fee                float64             `json:"fee"                exposure:"private,needPermission"`
	TotalProducts      int                 `json:"totalProducts"      exposure:"private,needPermission"`
}

type OrderShipmentWithAddress struct {
	OrderShipment
	ReceiverAddress UserAddress `json:"receiverAddress" exposure:"private,needPermission"`
}

type OrderWithFullInfo struct {
	OrderBase
	Payment       OrderPayment             `json:"payment"       exposure:"private,needPermission"`
	Shipment      OrderShipmentWithAddress `json:"shipment"      exposure:"private,needPermission"`
	TotalProducts int                      `json:"totalProducts" exposure:"private,needPermission"`
}

type OrderProductSelectedAttribute struct {
	Label       string `json:"label"       exposure:"private,needPermission"`
	Value       string `json:"value"       exposure:"private,needPermission"`
	AttributeId int    `json:"attributeId" exposure:"private,needPermission"`
	OptionId    int    `json:"optionId"    exposure:"private,needPermission"`
}

type OrderProductVariant struct {
	Id            int     `json:"id"            exposure:"private,needPermission"`
	Quantity      int     `json:"quantity"      exposure:"private,needPermission"`
	VariantPrice  float64 `json:"variantPrice"  exposure:"private,needPermission"`
	ShippingPrice float64 `json:"shippingPrice" exposure:"private,needPermission"`
	OrderId       int     `json:"orderId"       exposure:"private,needPermission"`
	VariantId     int     `json:"variantId"     exposure:"private,needPermission"`
}

type OrderProductVariantInfo struct {
	OrderProductVariant
	SelectedVariant ProductVariantWithAttributeSet `json:"selectedVariant" exposure:"private,needPermission"`
	Product         Product                        `json:"product"         exposure:"private,needPermission"`
}

type OrderProductVariantAssignmentPayload struct {
	Quantity  int `json:"quantity"  validate:"required"`
	VariantId int `json:"variantId" validate:"required"`
}

type CreateOrderPayload struct {
	UserId            int                                    `json:"userId"`
	ArrivalDate       time.Time                              `json:"arrivalDate"       validate:"required"`
	ProductVariants   []OrderProductVariantAssignmentPayload `json:"productVariants"   validate:"required"`
	ReceiverAddressId int                                    `json:"receiverAddressId" validate:"required"`
}

type OrderSearchQuery struct {
	UserId            *int                `json:"userId"`
	StoreId           *int                `json:"storeId"`
	PaymentStatus     *OrderPaymentStatus `json:"paymentStatus"`
	ShipmentStatus    *OrderPaymentStatus `json:"shipmentStatus"`
	CreatedAtLessThan *time.Time          `json:"createdAtLessThan"`
	CreatedAtMoreThan *time.Time          `json:"createdAtMoreThan"`
	Limit             *int                `json:"limit"`
	Offset            *int                `json:"offset"`
}

type UpdateOrderShipmentPayload struct {
	Status      *OrderShipmentStatus `json:"status"`
	ArrivalDate *time.Time           `json:"arrivalDate"`
}

type UpdateOrderPaymentPayload struct {
	Status *OrderPaymentStatus `json:"status"`
}
