package types

import "time"

// OrderBase represents basic order information
// @model OrderBase
type OrderBase struct {
	// Unique identifier for the order (private, needs permission)
	Id int `json:"id"        exposure:"private,needPermission"`
	// When the order was created (private, needs permission)
	CreatedAt time.Time `json:"createdAt" exposure:"private,needPermission"`
	// When the order was last updated (private, needs permission)
	UpdatedAt time.Time `json:"updatedAt" exposure:"private,needPermission"`
	// ID of the user who placed the order (private, needs permission)
	UserId int `json:"userId"    exposure:"private,needPermission"`
}

// OrderPayment represents payment information for an order
// @model OrderPayment
type OrderPayment struct {
	// Unique identifier for the payment (private, needs permission)
	Id int `json:"id"                 exposure:"private,needPermission"`
	// Total price of all product variants in the order (private, needs permission)
	TotalVariantsPrice float64 `json:"totalVariantsPrice" exposure:"private,needPermission"`
	// Total shipping cost for the order (private, needs permission)
	TotalShipmentPrice float64 `json:"totalShipmentPrice" exposure:"private,needPermission"`
	// Any additional fees applied to the order (private, needs permission)
	Fee float64 `json:"fee"                exposure:"private,needPermission"`
	// Current status of the payment (private, needs permission)
	Status OrderPaymentStatus `json:"status"             exposure:"private,needPermission"`
	// When the payment was created (private, needs permission)
	CreatedAt time.Time `json:"createdAt"          exposure:"private,needPermission"`
	// When the payment was last updated (private, needs permission)
	UpdatedAt time.Time `json:"updatedAt"          exposure:"private,needPermission"`
	// ID of the order this payment belongs to (private, needs permission)
	OrderId int `json:"orderId"            exposure:"private,needPermission"`
}

// OrderShipment represents shipment information for an order
// @model OrderShipment
type OrderShipment struct {
	// Unique identifier for the shipment (private, needs permission)
	Id int `json:"id"                exposure:"private,needPermission"`
	// Estimated arrival date of the shipment (private, needs permission)
	ArrivalDate time.Time `json:"arrivalDate"       exposure:"private,needPermission"`
	// Current status of the shipment (private, needs permission)
	Status OrderShipmentStatus `json:"status"            exposure:"private,needPermission"`
	// When the shipment was created (private, needs permission)
	CreatedAt time.Time `json:"createdAt"         exposure:"private,needPermission"`
	// When the shipment was last updated (private, needs permission)
	UpdatedAt time.Time `json:"updatedAt"         exposure:"private,needPermission"`
	// ID of the order this shipment belongs to (private, needs permission)
	OrderId int `json:"orderId"           exposure:"private,needPermission"`
	// ID of the receiver's address (private, needs permission)
	ReceiverAddressId int `json:"receiverAddressId" exposure:"private,needPermission"`
}

// Order represents a complete order with summarized information
// @model Order
type Order struct {
	OrderBase
	// Current status of the payment (private, needs permission)
	PaymentStatus OrderPaymentStatus `json:"paymentStatus"      exposure:"private,needPermission"`
	// Current status of the shipment (private, needs permission)
	ShipmentStatus OrderShipmentStatus `json:"shipmentStatus"     exposure:"private,needPermission"`
	// Total price of all product variants (private, needs permission)
	TotalVariantsPrice float64 `json:"totalVariantsPrice" exposure:"private,needPermission"`
	// Total shipping cost (private, needs permission)
	TotalShipmentPrice float64 `json:"totalShipmentPrice" exposure:"private,needPermission"`
	// Any additional fees (private, needs permission)
	Fee float64 `json:"fee"                exposure:"private,needPermission"`
	// Total number of products in the order (private, needs permission)
	TotalProducts int `json:"totalProducts"      exposure:"private,needPermission"`
}

// OrderShipmentWithAddress represents shipment information with the receiver's address
// @model OrderShipmentWithAddress
type OrderShipmentWithAddress struct {
	OrderShipment
	// Complete receiver address information (private, needs permission)
	ReceiverAddress UserAddress `json:"receiverAddress" exposure:"private,needPermission"`
}

// OrderWithFullInfo represents an order with complete payment and shipment details
// @model OrderWithFullInfo
type OrderWithFullInfo struct {
	OrderBase
	// Detailed payment information (private, needs permission)
	Payment OrderPayment `json:"payment"       exposure:"private,needPermission"`
	// Detailed shipment information with address (private, needs permission)
	Shipment OrderShipmentWithAddress `json:"shipment"      exposure:"private,needPermission"`
	// Total number of products in the order (private, needs permission)
	TotalProducts int `json:"totalProducts" exposure:"private,needPermission"`
}

// OrderProductSelectedAttribute represents selected attributes for an ordered product variant
// @model OrderProductSelectedAttribute
type OrderProductSelectedAttribute struct {
	// Display label for the attribute (private, needs permission)
	Label string `json:"label"       exposure:"private,needPermission"`
	// Selected value for the attribute (private, needs permission)
	Value string `json:"value"       exposure:"private,needPermission"`
	// ID of the product attribute (private, needs permission)
	AttributeId int `json:"attributeId" exposure:"private,needPermission"`
	// ID of the selected option (private, needs permission)
	OptionId int `json:"optionId"    exposure:"private,needPermission"`
}

// OrderProductVariant represents a product variant in an order
// @model OrderProductVariant
type OrderProductVariant struct {
	// Unique identifier for the order variant (private, needs permission)
	Id int `json:"id"            exposure:"private,needPermission"`
	// Quantity ordered (private, needs permission)
	Quantity int `json:"quantity"      exposure:"private,needPermission"`
	// Price per unit of the variant (private, needs permission)
	VariantPrice float64 `json:"variantPrice"  exposure:"private,needPermission"`
	// Shipping cost for this variant (private, needs permission)
	ShippingPrice float64 `json:"shippingPrice" exposure:"private,needPermission"`
	// ID of the order this variant belongs to (private, needs permission)
	OrderId int `json:"orderId"       exposure:"private,needPermission"`
	// ID of the product variant (private, needs permission)
	VariantId int `json:"variantId"     exposure:"private,needPermission"`
}

// OrderProductVariantInfo represents detailed information about an ordered product variant
// @model OrderProductVariantInfo
type OrderProductVariantInfo struct {
	OrderProductVariant
	// Complete information about the selected variant (private, needs permission)
	SelectedVariant ProductVariantWithAttributeSet `json:"selectedVariant" exposure:"private,needPermission"`
	// Information about the base product (private, needs permission)
	Product Product `json:"product"         exposure:"private,needPermission"`
}

// OrderProductVariantAssignmentPayload contains data for assigning a product variant to an order
// @model OrderProductVariantAssignmentPayload
type OrderProductVariantAssignmentPayload struct {
	// Number of units to order (required)
	Quantity int `json:"quantity"  validate:"required"`
	// ID of the product variant to order (required)
	VariantId int `json:"variantId" validate:"required"`
}

// CreateOrderPayload contains data needed to create a new order
// @model CreateOrderPayload
type CreateOrderPayload struct {
	// ID of the user placing the order
	UserId int `json:"userId"`
	// Expected arrival date for the order (required)
	ArrivalDate time.Time `json:"arrivalDate"       validate:"required"`
	// List of product variants to order (required)
	ProductVariants []OrderProductVariantAssignmentPayload `json:"productVariants"   validate:"required"`
	// ID of the receiver's address (required)
	ReceiverAddressId int `json:"receiverAddressId" validate:"required"`
}

// OrderSearchQuery contains parameters for searching orders
// @model OrderSearchQuery
type OrderSearchQuery struct {
	// Filter by user ID
	UserId *int `json:"userId"`
	// Filter by store ID
	StoreId *int `json:"storeId"`
	// Filter by payment status
	PaymentStatus *OrderPaymentStatus `json:"paymentStatus"`
	// Filter by shipment status
	ShipmentStatus *OrderPaymentStatus `json:"shipmentStatus"`
	// Filter orders created before this date
	CreatedAtLessThan *time.Time `json:"createdAtLessThan"`
	// Filter orders created after this date
	CreatedAtMoreThan *time.Time `json:"createdAtMoreThan"`
	// Maximum number of results to return
	Limit *int `json:"limit"`
	// Number of results to skip
	Offset *int `json:"offset"`
}

// UpdateOrderShipmentPayload contains data for updating order shipment information
// @model UpdateOrderShipmentPayload
type UpdateOrderShipmentPayload struct {
	// New status for the shipment
	Status *OrderShipmentStatus `json:"status"`
	// Updated estimated arrival date
	ArrivalDate *time.Time `json:"arrivalDate"`
}

// UpdateOrderPaymentPayload contains data for updating order payment information
// @model UpdateOrderPaymentPayload
type UpdateOrderPaymentPayload struct {
	// New status for the payment
	Status *OrderPaymentStatus `json:"status"`
}

// OrderProductVariantInsertData contains data for inserting a product variant into an order
// @model OrderProductVariantInsertData
type OrderProductVariantInsertData struct {
	// Number of units ordered
	Quantity int
	// Price per unit at time of order
	VariantPrice float64
	// Shipping cost per unit
	ShippingPrice float64
	// ID of the product variant
	VariantId int
	// ID of the order
	OrderId int
}
