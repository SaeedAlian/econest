package types

// HTTPError represents an error message from the API
// @model HTTPError
type HTTPError struct {
	// Error message string
	Message string `json:"message"`
}

// NewWalletTransactionResponse contains the new transaction id
// @model NewWalletTransactionResponse
type NewWalletTransactionResponse struct {
	// New transaction id
	TxId int `json:"txId"`
}

// NewUserResponse contains the new user id
// @model NewUserResponse
type NewUserResponse struct {
	// New user id
	UserId int `json:"userId"`
}

// NewAddressResponse contains the new address id
// @model NewAddressResponse
type NewAddressResponse struct {
	// New address id
	AddressId int `json:"addressId"`
}

// NewPhoneNumberResponse contains the new phone number id
// @model NewPhoneNumberResponse
type NewPhoneNumberResponse struct {
	// New phone number id
	PhoneNumberId int `json:"phoneNumberId"`
}

// NewStoreResponse contains the new store id
// @model NewStoreResponse
type NewStoreResponse struct {
	// New store id
	StoreId int `json:"storeId"`
}

// NewRoleResponse contains the new role id
// @model NewRoleResponse
type NewRoleResponse struct {
	// New role id
	RoleId int `json:"roleId"`
}

// NewPermissionGroupResponse contains the new permission group id
// @model NewPermissionGroupResponse
type NewPermissionGroupResponse struct {
	// New permission group id
	PermissionGroupId int `json:"pgroupId"`
}

// NewProductResponse contains the new product id
// @model NewProductResponse
type NewProductResponse struct {
	// New product id
	ProductId int `json:"productId"`
}

// NewProductOfferResponse contains the new product offer id
// @model NewProductOfferResponse
type NewProductOfferResponse struct {
	// New product offer id
	OfferId int `json:"offerId"`
}

// NewProductAttributeResponse contains the new product attribute id
// @model NewProductAttributeResponse
type NewProductAttributeResponse struct {
	// New product attribute id
	AttributeId int `json:"attributeId"`
}

// NewProductCommentResponse contains the new product comment id
// @model NewProductCommentResponse
type NewProductCommentResponse struct {
	// New product comment id
	CommentId int `json:"commentId"`
}

// NewProductCategoryResponse contains the new product category id
// @model NewProductCategoryResponse
type NewProductCategoryResponse struct {
	// New product category id
	CategoryId int `json:"categoryId"`
}

// NewProductTagResponse contains the new product tag id
// @model NewProductTagResponse
type NewProductTagResponse struct {
	// New product tag id
	TagId int `json:"tagId"`
}

// NewOrderResponse contains the new order id
// @model NewOrderResponse
type NewOrderResponse struct {
	// New order id
	OrderId int `json:"orderId"`
}

// TotalPageCountResponse contains the total pages of a list
// @model TotalPageCountResponse
type TotalPageCountResponse struct {
	// Total pages count
	Pages int32 `json:"pages"`
}

// FileUploadResponse contains the uploaded file name
// @model FileUploadResponse
type FileUploadResponse struct {
	// Uploaded file name
	FileName string `json:"fileName"`
}

// ProductInventoryResponse contains the inventory information of product
// @model ProductInventoryResponse
type ProductInventoryResponse struct {
	// Total quantity of the product
	Total int `json:"total"`
	// If the quantity is greater than 0
	InStock bool `json:"inStock"`
}
