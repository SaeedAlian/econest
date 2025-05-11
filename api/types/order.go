package types

import (
	"database/sql"
	"time"
)

type Order struct {
	Id            int           `json:"id"            exposure:"private,needPermission"`
	Verified      bool          `json:"verified"      exposure:"private,needPermission"`
	Status        OrderStatus   `json:"status"        exposure:"private,needPermission"`
	CreatedAt     time.Time     `json:"createdAt"     exposure:"private,needPermission"`
	UpdatedAt     time.Time     `json:"updatedAt"     exposure:"private,needPermission"`
	UserId        int           `json:"userId"        exposure:"private,needPermission"`
	TransactionId sql.NullInt32 `json:"transactionId" exposure:"private,needPermission"`
}

type OrderShipment struct {
	Id                int          `json:"id"                exposure:"private,needPermission"`
	ArrivalDate       time.Time    `json:"arrivalDate"       exposure:"private,needPermission"`
	ShipmentDate      time.Time    `json:"shipmentDate"      exposure:"private,needPermission"`
	Status            OrderStatus  `json:"status"            exposure:"private,needPermission"`
	ShipmentType      ShipmentType `json:"shipmentType"      exposure:"private,needPermission"`
	CreatedAt         time.Time    `json:"createdAt"         exposure:"private,needPermission"`
	UpdatedAt         time.Time    `json:"updatedAt"         exposure:"private,needPermission"`
	OrderId           int          `json:"orderId"           exposure:"private,needPermission"`
	ReceiverAddressId int          `json:"receiverAddressId" exposure:"private,needPermission"`
	SenderAddressId   int          `json:"senderAddressId"   exposure:"private,needPermission"`
}

type OrderProductInfo struct {
	Id            int           `json:"id"                  exposure:"private,needPermission"`
	Name          string        `json:"name"                exposure:"private,needPermission"`
	Slug          string        `json:"slug"                exposure:"private,needPermission"`
	Price         float64       `json:"price"               exposure:"private,needPermission"`
	Description   string        `json:"description"         exposure:"private,needPermission"`
	CreatedAt     time.Time     `json:"createdAt"           exposure:"private,needPermission"`
	UpdatedAt     time.Time     `json:"updatedAt"           exposure:"private,needPermission"`
	SubcategoryId int           `json:"subcategoryId"       exposure:"private,needPermission"`
	TotalQuantity int           `json:"totalQuantity"       exposure:"private,needPermission"`
	Offer         *ProductOffer `json:"offer,omitempty"     exposure:"private,needPermission"`
	MainImage     *ProductImage `json:"mainImage,omitempty" exposure:"private,needPermission"`
}

type OrderProductSelectedAttribute struct {
	Label       string `json:"label"       exposure:"private,needPermission"`
	Value       string `json:"value"       exposure:"private,needPermission"`
	AttributeId int    `json:"attributeId" exposure:"private,needPermission"`
	OptionId    int    `json:"optionId"    exposure:"private,needPermission"`
}

type OrderProductVariant struct {
	Id        int `json:"id"        exposure:"private,needPermission"`
	Quantity  int `json:"quantity"  exposure:"private,needPermission"`
	OrderId   int `json:"orderId"   exposure:"private,needPermission"`
	VariantId int `json:"variantId" exposure:"private,needPermission"`
}

type OrderProductVariantInfo struct {
	Id                 int                             `json:"id"                 exposure:"private,needPermission"`
	Quantity           int                             `json:"quantity"           exposure:"private,needPermission"`
	VariantQuantity    int                             `json:"variantQuantity"    exposure:"private,needPermission"`
	SelectedAttributes []OrderProductSelectedAttribute `json:"selectedAttributes" exposure:"private,needPermission"`
	Product            ProductWithMainInfo             `json:"product"            exposure:"private,needPermission"`
}

type OrderProductVariantAssignmentPayload struct {
	Quantity  int `json:"quantity"  validate:"required"`
	VariantId int `json:"variantId" validate:"required"`
}

type CreateOrderPayload struct {
	UserId          int                                    `json:"userId"          validate:"required"`
	TransactionId   int                                    `json:"transactionId"   validate:"required"`
	ProductVariants []OrderProductVariantAssignmentPayload `json:"productVariants"`
}

type UpdateOrderPayload struct {
	Verified *bool        `json:"verified"`
	Status   *OrderStatus `json:"status"`
}

type OrderSearchQuery struct {
	UserId            *int         `json:"userId"`
	Verified          *bool        `json:"verified"`
	Status            *OrderStatus `json:"status"`
	CreatedAtLessThan *time.Time   `json:"createdAtLessThan"`
	CreatedAtMoreThan *time.Time   `json:"createdAtMoreThan"`
	Limit             *int         `json:"limit"`
	Offset            *int         `json:"offset"`
}

type CreateOrderShipmentPayload struct {
	ArrivalDate       time.Time    `json:"arrivalDate"       validate:"required"`
	ShipmentDate      time.Time    `json:"shipmentDate"      validate:"required"`
	ShipmentType      ShipmentType `json:"shipmentType"      validate:"required"`
	OrderId           int          `json:"orderId"           validate:"required"`
	ReceiverAddressId int          `json:"receiverAddressId" validate:"required"`
	SenderAddressId   int          `json:"senderAddressId"   validate:"required"`
}

type UpdateOrderShipmentPayload struct {
	Status            *OrderStatus `json:"status"`
	ArrivalDate       *time.Time   `json:"arrivalDate"`
	ReceiverAddressId *int         `json:"receiverAddressId"`
	SenderAddressId   *int         `json:"senderAddressId"`
}
