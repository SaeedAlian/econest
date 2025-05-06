package types

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
	for _, v := range ValidSettingVisibilityStatuses {
		if s == v {
			return true
		}
	}

	return false
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
	for _, v := range ValidCredentialVerificationStatuses {
		if s == v {
			return true
		}
	}

	return false
}

type TransactionType string

const (
	TransactionTypeDeposit  TransactionType = "deposit"
	TransactionTypeWithdraw TransactionType = "withdraw"
	TransactionTypePurchase TransactionType = "purchase"
)

var ValidTransactionTypes = []TransactionType{
	TransactionTypeDeposit,
	TransactionTypeWithdraw,
	TransactionTypePurchase,
}

func (t TransactionType) IsValid() bool {
	for _, v := range ValidTransactionTypes {
		if t == v {
			return true
		}
	}

	return false
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
	for _, v := range ValidTransactionStatuses {
		if s == v {
			return true
		}
	}

	return false
}

type Action string

const (
	ActionFullControl                      Action = "full_control"
	ActionCanBanUser                       Action = "can_ban_user"
	ActionCanUnbanUser                     Action = "can_unban_user"
	ActionCanAddStore                      Action = "can_add_store"
	ActionCanUpdateStore                   Action = "can_update_store"
	ActionCanDeleteStore                   Action = "can_delete_store"
	ActionCanCreateOrder                   Action = "can_create_order"
	ActionCanModifyOrder                   Action = "can_modify_order"
	ActionCanCreateOrderShipment           Action = "can_create_order_shipment"
	ActionCanModifyOrderShipment           Action = "can_modify_order_shipment"
	ActionCanDefineProductOffer            Action = "can_define_product_offer"
	ActionCanModifyProductOffer            Action = "can_modify_product_offer"
	ActionCanCancelProductOffer            Action = "can_cancel_product_offer"
	ActionCanAddRole                       Action = "can_add_role"
	ActionCanDeleteRole                    Action = "can_delete_role"
	ActionCanModifyRole                    Action = "can_modify_role"
	ActionCanAddUserWithRole               Action = "can_add_user_with_role"
	ActionCanAddPermissionGroup            Action = "can_add_permission_group"
	ActionCanDeletePermissionGroup         Action = "can_delete_permission_group"
	ActionCanModifyPermissionGroup         Action = "can_modify_permission_group"
	ActionCanAssignPermissionGroupToRole   Action = "can_assign_permission_group_to_role"
	ActionCanRemovePermissionGroupFromRole Action = "can_remove_permission_group_from_role"
	ActionCanAddProductCategory            Action = "can_add_product_category"
	ActionCanModifyProductCategory         Action = "can_modify_product_category"
	ActionCanDeleteProductCategory         Action = "can_delete_product_category"
	ActionCanAddProduct                    Action = "can_add_product"
	ActionCanUpdateProduct                 Action = "can_update_product"
	ActionCanDeleteProduct                 Action = "can_delete_product"
	ActionCanCreateProductComment          Action = "can_create_product_comment"
	ActionCanUpdateProductComment          Action = "can_update_product_comment"
	ActionCanDeleteProductComment          Action = "can_delete_product_comment"
)

var ValidActions = []Action{
	ActionFullControl,
	ActionCanBanUser,
	ActionCanUnbanUser,
	ActionCanAddStore,
	ActionCanUpdateStore,
	ActionCanDeleteStore,
	ActionCanCreateOrder,
	ActionCanModifyOrder,
	ActionCanCreateOrderShipment,
	ActionCanModifyOrderShipment,
	ActionCanDefineProductOffer,
	ActionCanModifyProductOffer,
	ActionCanCancelProductOffer,
	ActionCanAddRole,
	ActionCanDeleteRole,
	ActionCanModifyRole,
	ActionCanAddUserWithRole,
	ActionCanAddPermissionGroup,
	ActionCanDeletePermissionGroup,
	ActionCanModifyPermissionGroup,
	ActionCanAssignPermissionGroupToRole,
	ActionCanRemovePermissionGroupFromRole,
	ActionCanAddProductCategory,
	ActionCanModifyProductCategory,
	ActionCanDeleteProductCategory,
	ActionCanAddProduct,
	ActionCanUpdateProduct,
	ActionCanDeleteProduct,
	ActionCanCreateProductComment,
	ActionCanUpdateProductComment,
	ActionCanDeleteProductComment,
}

func (a Action) IsValid() bool {
	for _, v := range ValidActions {
		if a == v {
			return true
		}
	}
	return false
}

type Resource string

const (
	ResourceFullAccess                   Resource = "full_access"
	ResourceRolesAndPermissions          Resource = "roles_and_permissions"
	ResourceUsersFullAccess              Resource = "users_full_access"
	ResourceUsersPublicOnly              Resource = "users_public_only"
	ResourceWalletTransactionsFullAccess Resource = "wallet_transactions_full_access"
	ResourceStoresFullAccess             Resource = "stores_full_access"
	ResourceStoresPublicOnly             Resource = "stores_public_only"
	ResourceOrdersFullAccess             Resource = "orders_full_access"
)

var ValidResources = []Resource{
	ResourceFullAccess,
	ResourceRolesAndPermissions,
	ResourceUsersFullAccess,
	ResourceUsersPublicOnly,
	ResourceWalletTransactionsFullAccess,
	ResourceStoresFullAccess,
	ResourceStoresPublicOnly,
	ResourceOrdersFullAccess,
}

func (r Resource) IsValid() bool {
	for _, v := range ValidResources {
		if r == v {
			return true
		}
	}

	return false
}

type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "pending_payment"
	OrderStatusPaymentPaid    OrderStatus = "payment_paid"
	OrderStatusCancelled      OrderStatus = "cancelled"
)

var ValidOrderStatuses = []OrderStatus{
	OrderStatusPendingPayment,
	OrderStatusPaymentPaid,
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

type ShipmentType string

const (
	ShipmentTypeShipping  ShipmentType = "shipping"
	ShipmentTypeReturning ShipmentType = "returning"
)

var ValidShipmentTypes = []ShipmentType{
	ShipmentTypeReturning,
	ShipmentTypeShipping,
}

func (t ShipmentType) IsValid() bool {
	for _, v := range ValidShipmentTypes {
		if t == v {
			return true
		}
	}

	return false
}

type ShipmentStatus string

const (
	ShipmentStatusOnTheWay  ShipmentStatus = "on_the_way"
	ShipmentStatusDelivered ShipmentStatus = "delivered"
	ShipmentStatusCancelled ShipmentStatus = "cancelled"
)

var ValidShipmentStatuses = []ShipmentStatus{
	ShipmentStatusOnTheWay,
	ShipmentStatusDelivered,
	ShipmentStatusCancelled,
}

func (s ShipmentStatus) IsValid() bool {
	for _, v := range ValidShipmentStatuses {
		if s == v {
			return true
		}
	}

	return false
}
