package permission_db_manager

import "time"

type Action string

const (
	ActionFullControl                      Action = "full_control"
	ActionCanBanUser                       Action = "can_ban_user"
	ActionCanUnbanUser                     Action = "can_unban_user"
	ActionCanAddStore                      Action = "can_add_store"
	ActionCanUpdateStore                   Action = "can_update_store"
	ActionCanDeleteStore                   Action = "can_delete_store"
	ActionCanCreateOrder                   Action = "can_create_order"
	ActionCanChangeOrderStatus             Action = "can_change_order_status"
	ActionCanVerifyOrder                   Action = "can_verify_order"
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
	ActionCanChangeOrderStatus,
	ActionCanVerifyOrder,
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

type PermissionGroup struct {
	Id          int
	Name        string
	Description string
	CreatedAt   time.Time
}

type RolePermissionGroup struct {
	RoleId            int
	PermissionGroupId int
}

type ResourcePermission struct {
	Id        int
	Resource  Resource
	CreatedAt time.Time
	GroupId   int
}

type ActionPermission struct {
	Id        int
	Action    Action
	CreatedAt time.Time
	GroupId   int
}
