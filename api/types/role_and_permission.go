package types

import (
	"time"

	json_types "github.com/SaeedAlian/econest/api/types/json"
)

type Role struct {
	Id          int                       `json:"id"          exposure:"needPermission"`
	Name        string                    `json:"name"        exposure:"needPermission"`
	Description json_types.JSONNullString `json:"description" exposure:"needPermission"`
	CreatedAt   time.Time                 `json:"createdAt"   exposure:"needPermission"`
	UpdatedAt   time.Time                 `json:"updatedAt"   exposure:"needPermission"`
}

type RoleWithPermissionGroups struct {
	Role
	PermissionGroups []PermissionGroup `json:"permissionGroups" exposure:"needPermission"`
}

type PermissionGroup struct {
	Id          int                       `json:"id"          exposure:"needPermission"`
	Name        string                    `json:"name"        exposure:"needPermission"`
	Description json_types.JSONNullString `json:"description" exposure:"needPermission"`
	CreatedAt   time.Time                 `json:"createdAt"   exposure:"needPermission"`
}

type PermissionGroupWithPermissions struct {
	PermissionGroup
	ResourcePermissions []GroupResourcePermissionInfo `json:"resourcePermissions" exposure:"needPermission"`
	ActionPermissions   []GroupActionPermissionInfo   `json:"actionPermissions"   exposure:"needPermission"`
}

type RoleGroupAssignment struct {
	RoleId            int `json:"roleId"            exposure:"needPermission"`
	PermissionGroupId int `json:"permissionGroupId" exposure:"needPermission"`
}

type GroupResourcePermission struct {
	CreatedAt time.Time `json:"createdAt" exposure:"needPermission"`
	Resource  Resource  `json:"resource"  exposure:"needPermission"`
	GroupId   int       `json:"groupId"   exposure:"needPermission"`
}

type GroupActionPermission struct {
	CreatedAt time.Time `json:"createdAt" exposure:"needPermission"`
	Action    Action    `json:"action"    exposure:"needPermission"`
	GroupId   int       `json:"groupId"   exposure:"needPermission"`
}

type GroupResourcePermissionInfo struct {
	Resource Resource `json:"resource" exposure:"needPermission"`
}

type GroupActionPermissionInfo struct {
	Action Action `json:"action" exposure:"needPermission"`
}

type CreateRolePayload struct {
	Name        string `json:"name"        validate:"required"`
	Description string `json:"description"`
}

type UpdateRolePayload struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type RolesSearchQuery struct {
	Name *string `json:"name"`
}

type CreatePermissionGroupPayload struct {
	Name        string `json:"name"        validate:"required"`
	Description string `json:"description"`
}

type UpdatePermissionGroupPayload struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type PermissionGroupSearchQuery struct {
	Name *string `json:"name"`
}

type RoleGroupAssignmentPayload struct {
	RoleId   int   `json:"roleId"`
	GroupIds []int `json:"groupIds"`
}

type GroupResourcePermissionAssignmentPayload struct {
	GroupId   int      `json:"groupId"`
	Resources []string `json:"resources"`
}

type GroupActionPermissionAssignmentPayload struct {
	GroupId int      `json:"groupId"`
	Actions []string `json:"actions"`
}
