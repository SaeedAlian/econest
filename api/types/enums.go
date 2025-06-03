package types

import "slices"

// SettingVisibilityStatus defines visibility levels for user settings
// @model SettingVisibilityStatus
type SettingVisibilityStatus string

const (
	// Setting is private (only visible to owner)
	SettingVisibilityStatusPrivate SettingVisibilityStatus = "private"
	// Setting is public (visible to everyone)
	SettingVisibilityStatusPublic SettingVisibilityStatus = "public"
	// Include both private and public settings
	SettingVisibilityStatusBoth SettingVisibilityStatus = "both"
)

var ValidSettingVisibilityStatuses = []SettingVisibilityStatus{
	SettingVisibilityStatusBoth,
	SettingVisibilityStatusPrivate,
	SettingVisibilityStatusPublic,
}

func (s SettingVisibilityStatus) IsValid() bool {
	return slices.Contains(ValidSettingVisibilityStatuses, s)
}

func (s SettingVisibilityStatus) String() string {
	return string(s)
}

// CredentialVerificationStatus defines verification states for user credentials
// @model CredentialVerificationStatus
type CredentialVerificationStatus string

const (
	// Credential has been verified
	CredentialVerificationStatusVerified CredentialVerificationStatus = "verified"
	// Credential has not been verified
	CredentialVerificationStatusNotVerified CredentialVerificationStatus = "not_verified"
	// Include both verified and unverified credentials
	CredentialVerificationStatusBoth CredentialVerificationStatus = "both"
)

var ValidCredentialVerificationStatuses = []CredentialVerificationStatus{
	CredentialVerificationStatusVerified,
	CredentialVerificationStatusNotVerified,
	CredentialVerificationStatusBoth,
}

func (s CredentialVerificationStatus) IsValid() bool {
	return slices.Contains(ValidCredentialVerificationStatuses, s)
}

func (s CredentialVerificationStatus) String() string {
	return string(s)
}

// TransactionType defines types of wallet transactions
// @model TransactionType
type TransactionType string

const (
	// Money being added to the wallet
	TransactionTypeDeposit TransactionType = "deposit"
	// Money being withdrawn from the wallet
	TransactionTypeWithdraw TransactionType = "withdraw"
)

var ValidTransactionTypes = []TransactionType{
	TransactionTypeDeposit,
	TransactionTypeWithdraw,
}

func (t TransactionType) IsValid() bool {
	return slices.Contains(ValidTransactionTypes, t)
}

func (t TransactionType) String() string {
	return string(t)
}

// TransactionStatus defines possible states of a transaction
// @model TransactionStatus
type TransactionStatus string

const (
	// Transaction is pending processing
	TransactionStatusPending TransactionStatus = "pending"
	// Transaction completed successfully
	TransactionStatusSuccessful TransactionStatus = "successful"
	// Transaction failed to complete
	TransactionStatusFailed TransactionStatus = "failed"
)

var ValidTransactionStatuses = []TransactionStatus{
	TransactionStatusPending,
	TransactionStatusSuccessful,
	TransactionStatusFailed,
}

func (s TransactionStatus) IsValid() bool {
	return slices.Contains(ValidTransactionStatuses, s)
}

func (s TransactionStatus) String() string {
	return string(s)
}

// OrderPaymentStatus defines payment states for orders
// @model OrderPaymentStatus
type OrderPaymentStatus string

const (
	// Payment is pending processing
	OrderPaymentStatusPending OrderPaymentStatus = "pending"
	// Payment completed successfully
	OrderPaymentStatusSuccessful OrderPaymentStatus = "successful"
	// Payment failed to process
	OrderPaymentStatusFailed OrderPaymentStatus = "failed"
)

var ValidOrderPaymentStatuses = []OrderPaymentStatus{
	OrderPaymentStatusPending,
	OrderPaymentStatusSuccessful,
	OrderPaymentStatusFailed,
}

func (s OrderPaymentStatus) IsValid() bool {
	return slices.Contains(ValidOrderPaymentStatuses, s)
}

func (s OrderPaymentStatus) String() string {
	return string(s)
}

// Action defines all possible permission actions in the system
// @model Action
type Action string

const (
	// Grants all possible permissions
	ActionFullControl Action = "full_control"

	// Permission to ban users
	ActionCanBanUser Action = "can_ban_user"
	// Permission to unban users
	ActionCanUnbanUser Action = "can_unban_user"

	// Permission to add stores
	ActionCanAddStore Action = "can_add_store"
	// Permission to update stores
	ActionCanUpdateStore Action = "can_update_store"
	// Permission to delete stores
	ActionCanDeleteStore Action = "can_delete_store"

	// Permission to create orders
	ActionCanCreateOrder Action = "can_create_order"
	// Permission to delete orders
	ActionCanDeleteOrder Action = "can_delete_order"
	// Permission to update order shipments
	ActionCanUpdateOrderShipment Action = "can_update_order_shipment"
	// Permission to complete order payments
	ActionCanCompleteOrderPayment Action = "can_complete_order_payment"
	// Permission to cancel order payments
	ActionCanCancelOrderPayment Action = "can_cancel_order_payment"

	// Permission to add roles
	ActionCanAddRole Action = "can_add_role"
	// Permission to update roles
	ActionCanUpdateRole Action = "can_update_role"
	// Permission to delete roles
	ActionCanDeleteRole Action = "can_delete_role"

	// Permission to add users with specific roles
	ActionCanAddUserWithRole Action = "can_add_user_with_role"

	// Permission to add permission groups
	ActionCanAddPermissionGroup Action = "can_add_permission_group"
	// Permission to delete permission groups
	ActionCanDeletePermissionGroup Action = "can_delete_permission_group"
	// Permission to update permission groups
	ActionCanUpdatePermissionGroup Action = "can_update_permission_group"

	// Permission to assign permission groups to roles
	ActionCanAssignPermissionGroupToRole Action = "can_assign_permission_group_to_role"
	// Permission to remove permission groups from roles
	ActionCanRemovePermissionGroupFromRole Action = "can_remove_permission_group_from_role"

	// Permission to assign a resource or action permission to a group
	ActionCanAssignPermissionToGroup Action = "can_assign_permission_to_group"
	// Permission to remove a resource or action permission from a group
	ActionCanRemovePermissionFromGroup Action = "can_remove_permission_from_group"

	// Permission to add product categories
	ActionCanAddProductCategory Action = "can_add_product_category"
	// Permission to update product categories
	ActionCanUpdateProductCategory Action = "can_update_product_category"
	// Permission to delete product categories
	ActionCanDeleteProductCategory Action = "can_delete_product_category"

	// Permission to add products
	ActionCanAddProduct Action = "can_add_product"
	// Permission to update products
	ActionCanUpdateProduct Action = "can_update_product"
	// Permission to delete products
	ActionCanDeleteProduct Action = "can_delete_product"

	// Permission to add product tags
	ActionCanAddProductTag Action = "can_add_product_tag"
	// Permission to update product tags
	ActionCanUpdateProductTag Action = "can_update_product_tag"
	// Permission to delete product tags
	ActionCanDeleteProductTag Action = "can_delete_product_tag"

	// Permission to add product offers
	ActionCanAddProductOffer Action = "can_add_product_offer"
	// Permission to update product offers
	ActionCanUpdateProductOffer Action = "can_update_product_offer"
	// Permission to delete product offers
	ActionCanDeleteProductOffer Action = "can_delete_product_offer"

	// Permission to add product attributes
	ActionCanAddProductAttribute Action = "can_add_product_attribute"
	// Permission to update product attributes
	ActionCanUpdateProductAttribute Action = "can_update_product_attribute"
	// Permission to delete product attributes
	ActionCanDeleteProductAttribute Action = "can_delete_product_attribute"

	// Permission to delete product comments
	ActionCanDeleteProductComment Action = "can_delete_product_comment"

	// Permission to approve withdrawal transactions
	ActionCanApproveWithdrawTransaction Action = "can_approve_withdraw_transaction"
	// Permission to cancel withdrawal transactions
	ActionCanCancelWithdrawTransaction Action = "can_cancel_withdraw_transaction"
)

var ValidActions = []Action{
	ActionFullControl,

	ActionCanBanUser,
	ActionCanUnbanUser,

	ActionCanAddStore,
	ActionCanUpdateStore,
	ActionCanDeleteStore,

	ActionCanCreateOrder,
	ActionCanDeleteOrder,
	ActionCanUpdateOrderShipment,
	ActionCanCompleteOrderPayment,
	ActionCanCancelOrderPayment,

	ActionCanAddRole,
	ActionCanUpdateRole,
	ActionCanDeleteRole,

	ActionCanAddUserWithRole,

	ActionCanAddPermissionGroup,
	ActionCanDeletePermissionGroup,
	ActionCanUpdatePermissionGroup,

	ActionCanAssignPermissionGroupToRole,
	ActionCanRemovePermissionGroupFromRole,

	ActionCanAssignPermissionToGroup,
	ActionCanRemovePermissionFromGroup,

	ActionCanAddProductCategory,
	ActionCanUpdateProductCategory,
	ActionCanDeleteProductCategory,

	ActionCanAddProduct,
	ActionCanUpdateProduct,
	ActionCanDeleteProduct,

	ActionCanAddProductTag,
	ActionCanUpdateProductTag,
	ActionCanDeleteProductTag,

	ActionCanAddProductOffer,
	ActionCanUpdateProductOffer,
	ActionCanDeleteProductOffer,

	ActionCanAddProductAttribute,
	ActionCanUpdateProductAttribute,
	ActionCanDeleteProductAttribute,

	ActionCanDeleteProductComment,

	ActionCanApproveWithdrawTransaction,
	ActionCanCancelWithdrawTransaction,
}

func (a Action) IsValid() bool {
	return slices.Contains(ValidActions, a)
}

func (a Action) String() string {
	return string(a)
}

// Resource defines system resources that can be permission-controlled
// @model Resource
type Resource string

const (
	// Full access to all resources
	ResourceFullAccess Resource = "full_access"

	// Access to roles and permissions management
	ResourceRolesAndPermissions Resource = "roles_and_permissions"

	// Full access to user management
	ResourceUsersFullAccess Resource = "users_full_access"

	// Full access to wallet transactions
	ResourceWalletTransactionsFullAccess Resource = "wallet_transactions_full_access"

	// Full access to store management
	ResourceStoresFullAccess Resource = "stores_full_access"

	// Full access to order management
	ResourceOrdersFullAccess Resource = "orders_full_access"
)

var ValidResources = []Resource{
	ResourceFullAccess,
	ResourceRolesAndPermissions,
	ResourceUsersFullAccess,
	ResourceWalletTransactionsFullAccess,
	ResourceStoresFullAccess,
	ResourceOrdersFullAccess,
}

func (r Resource) IsValid() bool {
	return slices.Contains(ValidResources, r)
}

func (r Resource) String() string {
	return string(r)
}

// OrderShipmentStatus defines possible states of order shipments
// @model OrderShipmentStatus
type OrderShipmentStatus string

const (
	// Shipment details to be determined
	OrderShipmentStatusToBeDetermined OrderShipmentStatus = "to_be_determined"
	// Shipment is in transit
	OrderShipmentStatusOnTheWay OrderShipmentStatus = "on_the_way"
	// Shipment has been delivered
	OrderShipmentStatusDelivered OrderShipmentStatus = "delivered"
	// Shipment was cancelled
	OrderShipmentStatusCancelled OrderShipmentStatus = "cancelled"
)

var ValidOrderShipmentStatuses = []OrderShipmentStatus{
	OrderShipmentStatusToBeDetermined,
	OrderShipmentStatusOnTheWay,
	OrderShipmentStatusDelivered,
	OrderShipmentStatusCancelled,
}

func (s OrderShipmentStatus) IsValid() bool {
	return slices.Contains(ValidOrderShipmentStatuses, s)
}

func (s OrderShipmentStatus) String() string {
	return string(s)
}

// DefaultRole defines system default role types
// @model DefaultRole
type DefaultRole string

const (
	// Super Administrator role with full system access
	DefaultRoleSuperAdmin DefaultRole = "Super Admin"
	// Administrator role with elevated privileges
	DefaultRoleAdmin DefaultRole = "Admin"
	// Vendor role for store owners
	DefaultRoleVendor DefaultRole = "Vendor"
	// Customer role for regular users
	DefaultRoleCustomer DefaultRole = "Customer"
)

var ValidDefaultRoles = []DefaultRole{
	DefaultRoleSuperAdmin,
	DefaultRoleAdmin,
	DefaultRoleVendor,
	DefaultRoleCustomer,
}

func (r DefaultRole) IsValid() bool {
	return slices.Contains(ValidDefaultRoles, r)
}

func (r DefaultRole) String() string {
	return string(r)
}
