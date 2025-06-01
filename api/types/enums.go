package types

import "slices"

type SettingVisibilityStatus string

const (
	SettingVisibilityStatusPrivate SettingVisibilityStatus = "private"
	SettingVisibilityStatusPublic  SettingVisibilityStatus = "public"
	SettingVisibilityStatusBoth    SettingVisibilityStatus = "both"
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

type CredentialVerificationStatus string

const (
	CredentialVerificationStatusVerified    CredentialVerificationStatus = "verified"
	CredentialVerificationStatusNotVerified CredentialVerificationStatus = "not_verified"
	CredentialVerificationStatusBoth        CredentialVerificationStatus = "both"
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

type TransactionType string

const (
	TransactionTypeDeposit  TransactionType = "deposit"
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

type TransactionStatus string

const (
	TransactionStatusPending    TransactionStatus = "pending"
	TransactionStatusSuccessful TransactionStatus = "successful"
	TransactionStatusFailed     TransactionStatus = "failed"
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

type OrderPaymentStatus string

const (
	OrderPaymentStatusPending    OrderPaymentStatus = "pending"
	OrderPaymentStatusSuccessful OrderPaymentStatus = "successful"
	OrderPaymentStatusFailed     OrderPaymentStatus = "failed"
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

type Action string

const (
	ActionFullControl Action = "full_control"

	ActionCanBanUser   Action = "can_ban_user"
	ActionCanUnbanUser Action = "can_unban_user"

	ActionCanAddStore    Action = "can_add_store"
	ActionCanUpdateStore Action = "can_update_store"
	ActionCanDeleteStore Action = "can_delete_store"

	ActionCanCreateOrder Action = "can_create_order"
	ActionCanUpdateOrder Action = "can_update_order"

	ActionCanCreateOrderShipment Action = "can_create_order_shipment"
	ActionCanUpdateOrderShipment Action = "can_update_order_shipment"

	ActionCanAddRole    Action = "can_add_role"
	ActionCanUpdateRole Action = "can_update_role"
	ActionCanDeleteRole Action = "can_delete_role"

	ActionCanAddUserWithRole Action = "can_add_user_with_role"

	ActionCanAddPermissionGroup    Action = "can_add_permission_group"
	ActionCanDeletePermissionGroup Action = "can_delete_permission_group"
	ActionCanUpdatePermissionGroup Action = "can_update_permission_group"

	ActionCanAssignPermissionGroupToRole   Action = "can_assign_permission_group_to_role"
	ActionCanRemovePermissionGroupFromRole Action = "can_remove_permission_group_from_role"

	ActionCanAssignPermissionToGroup   Action = "can_assign_permission_to_group"
	ActionCanRemovePermissionFromGroup Action = "can_remove_permission_from_group"

	ActionCanAddProductCategory    Action = "can_add_product_category"
	ActionCanUpdateProductCategory Action = "can_update_product_category"
	ActionCanDeleteProductCategory Action = "can_delete_product_category"

	ActionCanAddProduct    Action = "can_add_product"
	ActionCanUpdateProduct Action = "can_update_product"
	ActionCanDeleteProduct Action = "can_delete_product"

	ActionCanAddProductTag    Action = "can_add_product_tag"
	ActionCanUpdateProductTag Action = "can_update_product_tag"
	ActionCanDeleteProductTag Action = "can_delete_product_tag"

	ActionCanAddProductOffer    Action = "can_add_product_offer"
	ActionCanUpdateProductOffer Action = "can_update_product_offer"
	ActionCanDeleteProductOffer Action = "can_delete_product_offer"

	ActionCanAddProductAttribute    Action = "can_add_product_attribute"
	ActionCanUpdateProductAttribute Action = "can_update_product_attribute"
	ActionCanDeleteProductAttribute Action = "can_delete_product_attribute"

	ActionCanDeleteProductComment Action = "can_delete_product_comment"
)

var ValidActions = []Action{
	ActionFullControl,

	ActionCanBanUser,
	ActionCanUnbanUser,

	ActionCanAddStore,
	ActionCanUpdateStore,
	ActionCanDeleteStore,

	ActionCanCreateOrder,
	ActionCanUpdateOrder,

	ActionCanCreateOrderShipment,
	ActionCanUpdateOrderShipment,

	ActionCanAddRole,
	ActionCanUpdateRole,
	ActionCanDeleteRole,

	ActionCanAddUserWithRole,

	ActionCanAddPermissionGroup,
	ActionCanDeletePermissionGroup,
	ActionCanUpdatePermissionGroup,

	ActionCanAssignPermissionGroupToRole,
	ActionCanRemovePermissionGroupFromRole,

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
}

func (a Action) IsValid() bool {
	return slices.Contains(ValidActions, a)
}

func (a Action) String() string {
	return string(a)
}

type Resource string

const (
	ResourceFullAccess Resource = "full_access"

	ResourceRolesAndPermissions Resource = "roles_and_permissions"

	ResourceUsersFullAccess Resource = "users_full_access"

	ResourceWalletTransactionsFullAccess Resource = "wallet_transactions_full_access"

	ResourceStoresFullAccess Resource = "stores_full_access"

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

type OrderShipmentStatus string

const (
	OrderShipmentStatusToBeDetermined OrderShipmentStatus = "to_be_determined"
	OrderShipmentStatusOnTheWay       OrderShipmentStatus = "on_the_way"
	OrderShipmentStatusDelivered      OrderShipmentStatus = "delivered"
	OrderShipmentStatusCancelled      OrderShipmentStatus = "cancelled"
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

type DefaultRole string

const (
	DefaultRoleSuperAdmin DefaultRole = "Super Admin"
	DefaultRoleAdmin      DefaultRole = "Admin"
	DefaultRoleVendor     DefaultRole = "Vendor"
	DefaultRoleCustomer   DefaultRole = "Customer"
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
