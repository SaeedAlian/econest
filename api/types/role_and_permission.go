package types

import "time"

type Role struct {
	Id          int       `json:"id"          exposure:"needPermission"`
	Name        string    `json:"name"        exposure:"needPermission"`
	Description string    `json:"description" exposure:"needPermission"`
	CreatedAt   time.Time `json:"createdAt"   exposure:"needPermission"`
	UpdatedAt   time.Time `json:"updatedAt"   exposure:"needPermission"`
}

type RoleWithPermissionGroups struct {
	Id               int               `json:"id"               exposure:"needPermission"`
	Name             string            `json:"name"             exposure:"needPermission"`
	Description      string            `json:"description"      exposure:"needPermission"`
	CreatedAt        time.Time         `json:"createdAt"        exposure:"needPermission"`
	UpdatedAt        time.Time         `json:"updatedAt"        exposure:"needPermission"`
	PermissionGroups []PermissionGroup `json:"permissionGroups" exposure:"needPermission"`
}

type PermissionGroup struct {
	Id          int       `json:"id"          exposure:"needPermission"`
	Name        string    `json:"name"        exposure:"needPermission"`
	Description string    `json:"description" exposure:"needPermission"`
	CreatedAt   time.Time `json:"createdAt"   exposure:"needPermission"`
}

type PermissionGroupWithPermissions struct {
	Id                  int                           `json:"id"                  exposure:"needPermission"`
	Name                string                        `json:"name"                exposure:"needPermission"`
	Description         string                        `json:"description"         exposure:"needPermission"`
	CreatedAt           time.Time                     `json:"createdAt"           exposure:"needPermission"`
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

type CreateGroupAssignmentPayload struct {
	RoleId  int `json:"roleId"`
	GroupId int `json:"groupId"`
}

type CreateGroupResourcePermissionPayload struct {
	GroupId  int    `json:"groupId"`
	Resource string `json:"resource"`
}

type CreateGroupActionPermissionPayload struct {
	GroupId int    `json:"groupId"`
	Action  string `json:"action"`
}
