package permission_store

import "time"

type Action string

const (
	ActionCanAddProduct            Action = "can_add_product"
	ActionCanUpdateProduct         Action = "can_update_product"
	ActionCanDeleteProduct         Action = "can_delete_product"
	ActionCanAddVendor             Action = "can_add_vendor"
	ActionCanUpdateVendor          Action = "can_update_vendor"
	ActionCanDeleteVendor          Action = "can_delete_vendor"
	ActionCanBanUser               Action = "can_ban_user"
	ActionCanUnbanUser             Action = "can_unban_user"
	ActionCanAddProductTag         Action = "can_add_product_tag"
	ActionCanDeleteProductTag      Action = "can_delete_product_tag"
	ActionCanAddProductCategory    Action = "can_add_product_category"
	ActionCanDeleteProductCategory Action = "can_delete_product_category"
	ActionCanDeleteProductComment  Action = "can_delete_product_comment"
	ActionCanAddRole               Action = "can_add_role"
	ActionCanDeleteRole            Action = "can_delete_role"
	ActionCanModifyRole            Action = "can_modify_role"
	ActionCanAddPermissionGroup    Action = "can_add_permission_group"
	ActionCanDeletePermissionGroup Action = "can_delete_permission_group"
	ActionCanModifyPermissionGroup Action = "can_modify_permission_group"
)

var ValidActions = []Action{
	ActionCanAddProduct,
	ActionCanUpdateProduct,
	ActionCanDeleteProduct,
	ActionCanAddVendor,
	ActionCanUpdateVendor,
	ActionCanDeleteVendor,
	ActionCanBanUser,
	ActionCanUnbanUser,
	ActionCanAddProductTag,
	ActionCanDeleteProductTag,
	ActionCanAddProductCategory,
	ActionCanDeleteProductCategory,
	ActionCanDeleteProductComment,
	ActionCanAddRole,
	ActionCanDeleteRole,
	ActionCanModifyRole,
	ActionCanAddPermissionGroup,
	ActionCanDeletePermissionGroup,
	ActionCanModifyPermissionGroup,
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
	ResourceProducts    Resource = "products"
	ResourceOrders      Resource = "orders"
	ResourceUsers       Resource = "users"
	ResourcePermissions Resource = "permissions"
)

var ValidResources = []Resource{
	ResourceProducts,
	ResourceOrders,
	ResourceUsers,
	ResourcePermissions,
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
	Id          int
	Resource    Resource
	Description string
	CreatedAt   time.Time
	GroupId     int
}

type ActionPermission struct {
	Id          int
	Action      Action
	Description string
	CreatedAt   time.Time
	GroupId     int
}
