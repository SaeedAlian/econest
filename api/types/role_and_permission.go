package types

import (
	"time"

	json_types "github.com/SaeedAlian/econest/api/types/json"
)

// Role represents a user role in the system
// @model Role
type Role struct {
	// Unique identifier for the role (needs permission)
	Id int `json:"id"          exposure:"needPermission"`
	// Name of the role (needs permission)
	Name string `json:"name"        exposure:"needPermission"`
	// Description of the role (needs permission)
	Description json_types.JSONNullString `json:"description" exposure:"needPermission" swaggertype:"string"`
	// When the role was created (needs permission)
	CreatedAt time.Time `json:"createdAt"   exposure:"needPermission"`
	// When the role was last updated (needs permission)
	UpdatedAt time.Time `json:"updatedAt"   exposure:"needPermission"`
}

// RoleWithPermissionGroups combines Role with its associated permission groups
// @model RoleWithPermissionGroups
type RoleWithPermissionGroups struct {
	Role
	// List of permission groups assigned to this role (needs permission)
	PermissionGroups []PermissionGroup `json:"permissionGroups" exposure:"needPermission"`
}

// PermissionGroup represents a group of permissions
// @model PermissionGroup
type PermissionGroup struct {
	// Unique identifier for the permission group (needs permission)
	Id int `json:"id"          exposure:"needPermission"`
	// Name of the permission group (needs permission)
	Name string `json:"name"        exposure:"needPermission"`
	// Description of the permission group (needs permission)
	Description json_types.JSONNullString `json:"description" exposure:"needPermission" swaggertype:"string"`
	// When the permission group was created (needs permission)
	CreatedAt time.Time `json:"createdAt"   exposure:"needPermission"`
}

// PermissionGroupWithPermissions combines PermissionGroup with its resource and action permissions
// @model PermissionGroupWithPermissions
type PermissionGroupWithPermissions struct {
	PermissionGroup
	// List of resource permissions in this group (needs permission)
	ResourcePermissions []GroupResourcePermissionInfo `json:"resourcePermissions" exposure:"needPermission"`
	// List of action permissions in this group (needs permission)
	ActionPermissions []GroupActionPermissionInfo `json:"actionPermissions"   exposure:"needPermission"`
}

// RoleGroupAssignment represents an assignment of a permission group to a role
// @model RoleGroupAssignment
type RoleGroupAssignment struct {
	// ID of the role being assigned (needs permission)
	RoleId int `json:"roleId"            exposure:"needPermission"`
	// ID of the permission group being assigned (needs permission)
	PermissionGroupId int `json:"permissionGroupId" exposure:"needPermission"`
}

// GroupResourcePermission represents a resource permission assigned to a group
// @model GroupResourcePermission
type GroupResourcePermission struct {
	// When the permission was assigned (needs permission)
	CreatedAt time.Time `json:"createdAt" exposure:"needPermission"`
	// Resource being permitted (needs permission)
	Resource Resource `json:"resource"  exposure:"needPermission"`
	// ID of the group this permission belongs to (needs permission)
	GroupId int `json:"groupId"   exposure:"needPermission"`
}

// GroupActionPermission represents an action permission assigned to a group
// @model GroupActionPermission
type GroupActionPermission struct {
	// When the permission was assigned (needs permission)
	CreatedAt time.Time `json:"createdAt" exposure:"needPermission"`
	// Action being permitted (needs permission)
	Action Action `json:"action"    exposure:"needPermission"`
	// ID of the group this permission belongs to (needs permission)
	GroupId int `json:"groupId"   exposure:"needPermission"`
}

// GroupResourcePermissionInfo contains information about a resource permission
// @model GroupResourcePermissionInfo
type GroupResourcePermissionInfo struct {
	// Resource being permitted (needs permission)
	Resource Resource `json:"resource" exposure:"needPermission"`
}

// GroupActionPermissionInfo contains information about an action permission
// @model GroupActionPermissionInfo
type GroupActionPermissionInfo struct {
	// Action being permitted (needs permission)
	Action Action `json:"action" exposure:"needPermission"`
}

// CreateRolePayload contains data needed to create a new role
// @model CreateRolePayload
type CreateRolePayload struct {
	// Name of the role (required, 1+ characters)
	Name string `json:"name"        validate:"required"`
	// Description of the role
	Description string `json:"description"`
}

// UpdateRolePayload contains data for updating a role
// @model UpdateRolePayload
type UpdateRolePayload struct {
	// New name for the role (1+ characters if provided)
	Name *string `json:"name"`
	// New description for the role
	Description *string `json:"description"`
}

// RolesSearchQuery contains parameters for searching roles
// @model RolesSearchQuery
type RolesSearchQuery struct {
	// Filter by role name (partial match)
	Name *string `json:"name"`
}

// CreatePermissionGroupPayload contains data needed to create a new permission group
// @model CreatePermissionGroupPayload
type CreatePermissionGroupPayload struct {
	// Name of the permission group (required, 1+ characters)
	Name string `json:"name"        validate:"required"`
	// Description of the permission group
	Description string `json:"description"`
}

// UpdatePermissionGroupPayload contains data for updating a permission group
// @model UpdatePermissionGroupPayload
type UpdatePermissionGroupPayload struct {
	// New name for the permission group (1+ characters if provided)
	Name *string `json:"name"`
	// New description for the permission group
	Description *string `json:"description"`
}

// PermissionGroupSearchQuery contains parameters for searching permission groups
// @model PermissionGroupSearchQuery
type PermissionGroupSearchQuery struct {
	// Filter by permission group name (partial match)
	Name *string `json:"name"`
}

// RoleGroupAssignmentPayload contains data for assigning permission groups to a role
// @model RoleGroupAssignmentPayload
type RoleGroupAssignmentPayload struct {
	// ID of the role to assign groups to (required)
	RoleId int `json:"roleId"`
	// IDs of permission groups to assign (at least one required)
	GroupIds []int `json:"groupIds"`
}

// GroupResourcePermissionAssignmentPayload contains data for assigning resource permissions to a group
// @model GroupResourcePermissionAssignmentPayload
type GroupResourcePermissionAssignmentPayload struct {
	// ID of the permission group (required)
	GroupId int `json:"groupId"`
	// Names of resources to assign (at least one required)
	Resources []string `json:"resources"`
}

// GroupActionPermissionAssignmentPayload contains data for assigning action permissions to a group
// @model GroupActionPermissionAssignmentPayload
type GroupActionPermissionAssignmentPayload struct {
	// ID of the permission group (required)
	GroupId int `json:"groupId"`
	// Names of actions to assign (at least one required)
	Actions []string `json:"actions"`
}
