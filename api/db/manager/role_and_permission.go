package db_manager

import (
	"database/sql"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/lib/pq"

	"github.com/SaeedAlian/econest/api/types"
)

func (m *Manager) CreateRole(p types.CreateRolePayload) (int, error) {
	rowId := -1
	err := m.db.QueryRow("INSERT INTO roles (name, description) VALUES ($1, $2) RETURNING id;", p.Name, p.Description).
		Scan(&rowId)
	if err != nil {
		return -1, err
	}

	return rowId, nil
}

func (m *Manager) CreatePermissionGroup(p types.CreatePermissionGroupPayload) (int, error) {
	rowId := -1
	err := m.db.QueryRow("INSERT INTO permission_groups (name, description) VALUES ($1, $2) RETURNING id;",
		p.Name, p.Description,
	).
		Scan(&rowId)
	if err != nil {
		return -1, err
	}

	return rowId, nil
}

func (m *Manager) GetRoles(query types.RolesSearchQuery) ([]types.Role, error) {
	var base string
	base = "SELECT * FROM roles"

	q, args := buildRoleSearchQuery(query, base)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := []types.Role{}

	for rows.Next() {
		role, err := scanRoleRow(rows)
		if err != nil {
			return nil, err
		}

		roles = append(roles, *role)
	}

	return roles, nil
}

func (m *Manager) GetRoleById(id int) (*types.Role, error) {
	rows, err := m.db.Query(
		"SELECT * FROM roles WHERE id = $1;",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	role := new(types.Role)
	role.Id = -1

	for rows.Next() {
		role, err = scanRoleRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if role.Id == -1 {
		return nil, types.ErrRoleNotFound
	}

	return role, nil
}

func (m *Manager) GetRoleByName(name string) (*types.Role, error) {
	rows, err := m.db.Query(
		"SELECT * FROM roles WHERE name = $1;",
		name,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	role := new(types.Role)
	role.Id = -1

	for rows.Next() {
		role, err = scanRoleRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if role.Id == -1 {
		return nil, types.ErrRoleNotFound
	}

	return role, nil
}

func (m *Manager) GetRoleWithPermissionGroupsById(id int) (*types.RoleWithPermissionGroups, error) {
	rows, err := m.db.Query(
		"SELECT * FROM roles WHERE id = $1;",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	role := new(types.RoleWithPermissionGroups)
	role.Id = -1

	for rows.Next() {
		r, err := scanRoleRow(rows)
		if err != nil {
			return nil, err
		}

		role.Role = *r
	}

	if role.Id == -1 {
		return nil, types.ErrRoleNotFound
	}

	groupRows, err := m.db.Query(`SELECT
      pg.id, pg.name, pg.description, pg.created_at FROM permission_groups pg
      JOIN role_group_assignments rga ON pg.id = rga.permission_group_id 
      WHERE rga.role_id = $1;
    `, role.Id)
	if err != nil {
		return nil, err
	}
	defer groupRows.Close()

	groups := []types.PermissionGroup{}

	for groupRows.Next() {
		group, err := scanPermissionGroupRow(groupRows)
		if err != nil {
			return nil, err
		}

		groups = append(groups, *group)
	}

	role.PermissionGroups = groups

	return role, nil
}

func (m *Manager) GetRoleWithPermissionGroupsByName(
	name string,
) (*types.RoleWithPermissionGroups, error) {
	rows, err := m.db.Query(
		"SELECT * FROM roles WHERE name = $1;",
		name,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	role := new(types.RoleWithPermissionGroups)
	role.Id = -1

	for rows.Next() {
		r, err := scanRoleRow(rows)
		if err != nil {
			return nil, err
		}

		role.Role = *r
	}

	if role.Id == -1 {
		return nil, types.ErrRoleNotFound
	}

	groupRows, err := m.db.Query(`SELECT
      pg.id, pg.name, pg.description, pg.created_at FROM permission_groups pg
      JOIN role_group_assignments rga ON pg.id = rga.permission_group_id 
      WHERE rga.role_id = $1;
    `, role.Id)
	if err != nil {
		return nil, err
	}
	defer groupRows.Close()

	groups := []types.PermissionGroup{}

	for groupRows.Next() {
		group, err := scanPermissionGroupRow(groupRows)
		if err != nil {
			return nil, err
		}

		groups = append(groups, *group)
	}

	role.PermissionGroups = groups

	return role, nil
}

func (m *Manager) GetPermissionGroups(
	query types.PermissionGroupSearchQuery,
) ([]types.PermissionGroup, error) {
	var base string
	base = "SELECT * FROM permission_groups"

	q, args := buildPermissionGroupSearchQuery(query, base)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := []types.PermissionGroup{}

	for rows.Next() {
		group, err := scanPermissionGroupRow(rows)
		if err != nil {
			return nil, err
		}

		groups = append(groups, *group)
	}

	return groups, nil
}

func (m *Manager) GetPermissionGroupById(id int) (*types.PermissionGroup, error) {
	rows, err := m.db.Query(
		"SELECT * FROM permission_groups WHERE id = $1;",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	group := new(types.PermissionGroup)
	group.Id = -1

	for rows.Next() {
		group, err = scanPermissionGroupRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if group.Id == -1 {
		return nil, types.ErrPermissionGroupNotFound
	}

	return group, nil
}

func (m *Manager) GetPermissionGroupByName(name string) (*types.PermissionGroup, error) {
	rows, err := m.db.Query(
		"SELECT * FROM permission_groups WHERE name = $1;",
		name,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	group := new(types.PermissionGroup)
	group.Id = -1

	for rows.Next() {
		group, err = scanPermissionGroupRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if group.Id == -1 {
		return nil, types.ErrPermissionGroupNotFound
	}

	return group, nil
}

func (m *Manager) GetPermissionGroupWithPermissionsById(
	id int,
) (*types.PermissionGroupWithPermissions, error) {
	rows, err := m.db.Query(
		"SELECT * FROM permission_groups WHERE id = $1;",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	group := new(types.PermissionGroupWithPermissions)
	group.Id = -1

	for rows.Next() {
		g, err := scanPermissionGroupRow(rows)
		if err != nil {
			return nil, err
		}

		group.PermissionGroup = *g
	}

	if group.Id == -1 {
		return nil, types.ErrPermissionGroupNotFound
	}

	actionRows, err := m.db.Query(`SELECT
      gap.action FROM group_action_permissions gap
      JOIN permission_groups pg ON pg.id = gap.group_id
      WHERE gap.group_id = $1;
    `, group.Id)
	if err != nil {
		return nil, err
	}
	defer actionRows.Close()

	resourceRows, err := m.db.Query(`SELECT
      grp.resource FROM group_resource_permissions grp
      JOIN permission_groups pg ON pg.id = grp.group_id
      WHERE grp.group_id = $1;
    `, group.Id)
	if err != nil {
		return nil, err
	}
	defer resourceRows.Close()

	actions := []types.GroupActionPermissionInfo{}
	resources := []types.GroupResourcePermissionInfo{}

	for actionRows.Next() {
		action, err := scanActionPermissionInfoRow(actionRows)
		if err != nil {
			return nil, err
		}

		actions = append(actions, *action)
	}

	for resourceRows.Next() {
		resource, err := scanResourcePermissionInfoRow(resourceRows)
		if err != nil {
			return nil, err
		}

		resources = append(resources, *resource)
	}

	group.ResourcePermissions = resources
	group.ActionPermissions = actions

	return group, nil
}

func (m *Manager) GetPermissionGroupWithPermissionsByName(
	name string,
) (*types.PermissionGroupWithPermissions, error) {
	rows, err := m.db.Query(
		"SELECT * FROM permission_groups WHERE name = $1;",
		name,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	group := new(types.PermissionGroupWithPermissions)
	group.Id = -1

	for rows.Next() {
		g, err := scanPermissionGroupRow(rows)
		if err != nil {
			return nil, err
		}

		group.PermissionGroup = *g
	}

	if group.Id == -1 {
		return nil, types.ErrPermissionGroupNotFound
	}

	actionRows, err := m.db.Query(`SELECT
      gap.action FROM group_action_permissions gap
      JOIN permission_groups pg ON pg.id = gap.group_id
      WHERE gap.group_id = $1;
    `, group.Id)
	if err != nil {
		return nil, err
	}
	defer actionRows.Close()

	resourceRows, err := m.db.Query(`SELECT
      grp.resource FROM group_resource_permissions grp
      JOIN permission_groups pg ON pg.id = grp.group_id
      WHERE grp.group_id = $1;
    `, group.Id)
	if err != nil {
		return nil, err
	}
	defer resourceRows.Close()

	actions := []types.GroupActionPermissionInfo{}
	resources := []types.GroupResourcePermissionInfo{}

	for actionRows.Next() {
		action, err := scanActionPermissionInfoRow(actionRows)
		if err != nil {
			return nil, err
		}

		actions = append(actions, *action)
	}

	for resourceRows.Next() {
		resource, err := scanResourcePermissionInfoRow(resourceRows)
		if err != nil {
			return nil, err
		}

		resources = append(resources, *resource)
	}

	group.ResourcePermissions = resources
	group.ActionPermissions = actions

	return group, nil
}

func (m *Manager) GetRolesWithPermissionGroups(
	query types.RolesSearchQuery,
) ([]types.RoleWithPermissionGroups, error) {
	var base string
	base = "SELECT * FROM roles"

	q, args := buildRoleSearchQuery(query, base)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := []types.RoleWithPermissionGroups{}

	for rows.Next() {
		role, err := scanRoleRow(rows)
		if err != nil {
			return nil, err
		}

		groupRows, err := m.db.Query(`SELECT
      pg.id, pg.name, pg.description, pg.created_at FROM permission_groups pg
      JOIN role_group_assignments rga ON pg.id = rga.permission_group_id 
      WHERE rga.role_id = $1;
    `, role.Id)
		if err != nil {
			return nil, err
		}
		defer groupRows.Close()

		groups := []types.PermissionGroup{}

		for groupRows.Next() {
			group, err := scanPermissionGroupRow(groupRows)
			if err != nil {
				return nil, err
			}

			groups = append(groups, *group)
		}

		role_with_groups := types.RoleWithPermissionGroups{
			Role:             *role,
			PermissionGroups: groups,
		}

		roles = append(roles, role_with_groups)
	}

	return roles, nil
}

func (m *Manager) GetPermissionGroupsWithPermissions(
	query types.PermissionGroupSearchQuery,
) ([]types.PermissionGroupWithPermissions, error) {
	var base string
	base = "SELECT * FROM permission_groups"

	q, args := buildPermissionGroupSearchQuery(query, base)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := []types.PermissionGroupWithPermissions{}

	for rows.Next() {
		group, err := scanPermissionGroupRow(rows)
		if err != nil {
			return nil, err
		}

		actionRows, err := m.db.Query(`SELECT
      gap.action FROM group_action_permissions gap
      JOIN permission_groups pg ON pg.id = gap.group_id
      WHERE gap.group_id = $1;
    `, group.Id)
		if err != nil {
			return nil, err
		}
		defer actionRows.Close()

		resourceRows, err := m.db.Query(`SELECT
      grp.resource FROM group_resource_permissions grp
      JOIN permission_groups pg ON pg.id = grp.group_id
      WHERE grp.group_id = $1;
    `, group.Id)
		if err != nil {
			return nil, err
		}
		defer resourceRows.Close()

		actions := []types.GroupActionPermissionInfo{}
		resources := []types.GroupResourcePermissionInfo{}

		for actionRows.Next() {
			action, err := scanActionPermissionInfoRow(actionRows)
			if err != nil {
				return nil, err
			}

			actions = append(actions, *action)
		}

		for resourceRows.Next() {
			resource, err := scanResourcePermissionInfoRow(resourceRows)
			if err != nil {
				return nil, err
			}

			resources = append(resources, *resource)
		}

		group_with_permissions := types.PermissionGroupWithPermissions{
			PermissionGroup:     *group,
			ResourcePermissions: resources,
			ActionPermissions:   actions,
		}

		groups = append(groups, group_with_permissions)
	}

	return groups, nil
}

func (m *Manager) GetRolesBasedOnResourcePermission(
	resources []types.Resource,
) ([]types.Role, error) {
	rows, err := m.db.Query(`SELECT
    r.* FROM group_resource_permissions grp 
    JOIN role_group_assignments rga ON rga.permission_group_id = grp.group_id
		JOIN roles r ON rga.role_id = r.id WHERE grp.resource = ANY($1::resources[]);
  `, pq.Array(resources))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := []types.Role{}

	for rows.Next() {
		role, err := scanRoleRow(rows)
		if err != nil {
			return nil, err
		}

		roles = append(roles, *role)
	}

	return roles, nil
}

func (m *Manager) GetRolesBasedOnActionPermission(actions []types.Action) ([]types.Role, error) {
	rows, err := m.db.Query(`SELECT
    r.* FROM group_action_permissions gap 
    JOIN role_group_assignments rga ON rga.permission_group_id = gap.group_id
		JOIN roles r ON rga.role_id = r.id WHERE gap.action = ANY($1::actions[]);
  `, pq.Array(actions))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := []types.Role{}

	for rows.Next() {
		role, err := scanRoleRow(rows)
		if err != nil {
			return nil, err
		}

		roles = append(roles, *role)
	}

	return roles, nil
}

func (m *Manager) GetPermissionGroupsBasedOnResourcePermission(
	resources []types.Resource,
) ([]types.PermissionGroup, error) {
	rows, err := m.db.Query(`SELECT
    pg.* FROM group_resource_permissions grp
    JOIN permission_groups pg ON pg.id = grp.group_id
    WHERE grp.resource = ANY($1::resources[]);
  `, pq.Array(resources))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := []types.PermissionGroup{}

	for rows.Next() {
		group, err := scanPermissionGroupRow(rows)
		if err != nil {
			return nil, err
		}

		groups = append(groups, *group)
	}

	return groups, nil
}

func (m *Manager) GetPermissionGroupsBasedOnActionPermission(
	actions []types.Action,
) ([]types.PermissionGroup, error) {
	rows, err := m.db.Query(`SELECT
    pg.* FROM group_action_permissions gap
    JOIN permission_groups pg ON pg.id = gap.group_id
    WHERE gap.action = ANY($1::actions[]);
  `, pq.Array(actions))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := []types.PermissionGroup{}

	for rows.Next() {
		group, err := scanPermissionGroupRow(rows)
		if err != nil {
			return nil, err
		}

		groups = append(groups, *group)
	}

	return groups, nil
}

func (m *Manager) AddPermissionGroupsToRole(roleId int, permissionGroupIds []int) error {
	permissionGroupIdsLen := len(permissionGroupIds)
	if permissionGroupIdsLen == 0 {
		return nil
	}

	valueSqls := make([]string, 0, permissionGroupIdsLen)
	valueArgs := make([]any, 0, permissionGroupIdsLen*2)

	for i, permissionGroupId := range permissionGroupIds {
		valueSqls = append(valueSqls, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
		valueArgs = append(valueArgs, roleId, permissionGroupId)
	}

	query := fmt.Sprintf(
		"INSERT INTO role_group_assignments (role_id, permission_group_id) VALUES %s",
		strings.Join(valueSqls, ", "),
	)

	_, err := m.db.Exec(query, valueArgs...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) RemovePermissionGroupsFromRole(roleId int, permissionGroupIds []int) error {
	permissionGroupIdsLen := len(permissionGroupIds)
	if permissionGroupIdsLen == 0 {
		return nil
	}

	valueArgs := make([]any, 0, permissionGroupIdsLen+1)
	placeholders := make([]string, permissionGroupIdsLen)

	for i, permissionGroupId := range permissionGroupIds {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		valueArgs = append(valueArgs, permissionGroupId)
	}

	query := fmt.Sprintf(
		"DELETE FROM role_group_assignments WHERE role_id = $1 AND permission_group_id IN (%s)",
		strings.Join(placeholders, ", "),
	)

	valueArgs = append([]any{roleId}, valueArgs...)

	_, err := m.db.Exec(query, valueArgs...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) AddResourcePermissionsToGroup(
	groupId int,
	resources []types.Resource,
) error {
	resourcesLen := len(resources)
	if resourcesLen == 0 {
		return nil
	}

	valueSqls := make([]string, 0, resourcesLen)
	valueArgs := make([]any, 0, resourcesLen*2)

	for i, resource := range resources {
		valueSqls = append(valueSqls, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
		valueArgs = append(valueArgs, groupId, resource)
	}

	query := fmt.Sprintf(
		"INSERT INTO group_resource_permissions (group_id, resource) VALUES %s",
		strings.Join(valueSqls, ", "),
	)

	_, err := m.db.Exec(query, valueArgs...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) AddActionPermissionsToGroup(
	groupId int,
	actions []types.Action,
) error {
	actionsLen := len(actions)
	if actionsLen == 0 {
		return nil
	}

	valueSqls := make([]string, 0, actionsLen)
	valueArgs := make([]any, 0, actionsLen*2)

	for i, action := range actions {
		valueSqls = append(valueSqls, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
		valueArgs = append(valueArgs, groupId, action)
	}

	query := fmt.Sprintf(
		"INSERT INTO group_action_permissions (group_id, action) VALUES %s",
		strings.Join(valueSqls, ", "),
	)

	_, err := m.db.Exec(query, valueArgs...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) RemoveResourcePermissionsFromGroup(
	groupId int,
	resources []types.Resource,
) error {
	resourcesLen := len(resources)
	if resourcesLen == 0 {
		return nil
	}

	valueArgs := make([]any, 0, resourcesLen+1)
	placeholders := make([]string, resourcesLen)

	for i, resource := range resources {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		valueArgs = append(valueArgs, resource)
	}

	query := fmt.Sprintf(
		"DELETE FROM group_resource_permissions WHERE group_id = $1 AND resource IN (%s)",
		strings.Join(placeholders, ", "),
	)

	valueArgs = append([]any{groupId}, valueArgs...)

	_, err := m.db.Exec(query, valueArgs...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) RemoveActionPermissionsFromGroup(
	groupId int,
	actions []types.Action,
) error {
	actionsLen := len(actions)
	if actionsLen == 0 {
		return nil
	}

	valueArgs := make([]any, 0, actionsLen+1)
	placeholders := make([]string, actionsLen)

	for i, action := range actions {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		valueArgs = append(valueArgs, action)
	}

	query := fmt.Sprintf(
		"DELETE FROM group_action_permissions WHERE group_id = $1 AND action IN (%s)",
		strings.Join(placeholders, ", "),
	)

	valueArgs = append([]any{groupId}, valueArgs...)

	_, err := m.db.Exec(query, valueArgs...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) UpdateRole(id int, p types.UpdateRolePayload) error {
	clauses := []string{}
	args := []any{}
	argsPos := 1

	if p.Name != nil {
		clauses = append(clauses, fmt.Sprintf("name = $%d", argsPos))
		args = append(args, *p.Name)
		argsPos++
	}

	if p.Description != nil {
		clauses = append(clauses, fmt.Sprintf("description = $%d", argsPos))
		args = append(args, *p.Description)
		argsPos++
	}

	if len(clauses) == 0 {
		return types.ErrNoFieldsReceivedToUpdate
	}

	clauses = append(clauses, fmt.Sprintf("updated_at = $%d", argsPos))
	args = append(args, time.Now())
	argsPos++

	args = append(args, id)
	q := fmt.Sprintf(
		"UPDATE roles SET %s WHERE id = $%d",
		strings.Join(clauses, ", "),
		argsPos,
	)

	_, err := m.db.Exec(q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) UpdatePermissionGroup(id int, p types.UpdatePermissionGroupPayload) error {
	clauses := []string{}
	args := []any{}
	argsPos := 1

	if p.Name != nil {
		clauses = append(clauses, fmt.Sprintf("name = $%d", argsPos))
		args = append(args, *p.Name)
		argsPos++
	}

	if p.Description != nil {
		clauses = append(clauses, fmt.Sprintf("description = $%d", argsPos))
		args = append(args, *p.Description)
		argsPos++
	}

	if len(clauses) == 0 {
		return types.ErrNoFieldsReceivedToUpdate
	}

	clauses = append(clauses, fmt.Sprintf("updated_at = $%d", argsPos))
	args = append(args, time.Now())
	argsPos++

	args = append(args, id)
	q := fmt.Sprintf(
		"UPDATE permission_groups SET %s WHERE id = $%d",
		strings.Join(clauses, ", "),
		argsPos,
	)

	_, err := m.db.Exec(q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) DeleteRole(id int) error {
	_, err := m.db.Exec(
		"DELETE FROM roles WHERE id = $1;",
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) DeletePermissionGroup(id int) error {
	_, err := m.db.Exec(
		"DELETE FROM permission_groups WHERE id = $1;",
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) IsRoleHasAllActionPermissions(
	actions []types.Action,
	roleId int,
) (bool, error) {
	actionRows, err := m.db.Query(`
		SELECT gap.action FROM role_group_assignments rga 
		JOIN group_action_permissions gap ON gap.group_id = rga.permission_group_id
		WHERE rga.role_id = $1;
	`, roleId)
	if err != nil {
		return false, err
	}
	defer actionRows.Close()

	resultActions := []types.Action{}

	for actionRows.Next() {
		var a types.Action
		err := actionRows.Scan(&a)
		if err != nil {
			return false, err
		}
		resultActions = append(resultActions, a)
	}

	found := 0

	for _, r1 := range actions {
		for _, r2 := range resultActions {
			if r1 == r2 {
				found++
			}
		}
	}

	return found == len(actions), nil
}

func (m *Manager) IsRoleHasSomeActionPermissions(
	actions []types.Action,
	roleId int,
) (bool, error) {
	actionRows, err := m.db.Query(`
		SELECT gap.action FROM role_group_assignments rga 
		JOIN group_action_permissions gap ON gap.group_id = rga.permission_group_id
		WHERE rga.role_id = $1;
	`, roleId)
	if err != nil {
		return false, err
	}
	defer actionRows.Close()

	resultActions := []types.Action{}

	for actionRows.Next() {
		var a types.Action
		err := actionRows.Scan(&a)
		if err != nil {
			return false, err
		}
		resultActions = append(resultActions, a)
	}

	for _, r := range actions {
		if slices.Contains(resultActions, r) {
			return true, nil
		}
	}

	return false, nil
}

func (m *Manager) IsRoleHasAllResourcePermissions(
	resources []types.Resource,
	roleId int,
) (bool, error) {
	resourceRows, err := m.db.Query(`
		SELECT grp.resource FROM role_group_assignments rga 
		JOIN group_resource_permissions grp ON grp.group_id = rga.permission_group_id
		WHERE rga.role_id = $1;
	`, roleId)
	if err != nil {
		return false, err
	}
	defer resourceRows.Close()

	resultResources := []types.Resource{}

	for resourceRows.Next() {
		var r types.Resource
		err := resourceRows.Scan(&r)
		if err != nil {
			return false, err
		}
		resultResources = append(resultResources, r)
	}

	found := 0

	for _, r1 := range resources {
		for _, r2 := range resultResources {
			if r1 == r2 {
				found++
			}
		}
	}

	return found == len(resources), nil
}

func (m *Manager) IsRoleHasSomeResourcePermissions(
	resources []types.Resource,
	roleId int,
) (bool, error) {
	resourceRows, err := m.db.Query(`
		SELECT grp.resource FROM role_group_assignments rga 
		JOIN group_resource_permissions grp ON grp.group_id = rga.permission_group_id
		WHERE rga.role_id = $1;
	`, roleId)
	if err != nil {
		return false, err
	}
	defer resourceRows.Close()

	resultResources := []types.Resource{}

	for resourceRows.Next() {
		var r types.Resource
		err := resourceRows.Scan(&r)
		if err != nil {
			return false, err
		}
		resultResources = append(resultResources, r)
	}

	for _, r := range resources {
		if slices.Contains(resultResources, r) {
			return true, nil
		}
	}

	return false, nil
}

func scanRoleRow(rows *sql.Rows) (*types.Role, error) {
	n := new(types.Role)

	err := rows.Scan(
		&n.Id,
		&n.Name,
		&n.Description,
		&n.CreatedAt,
		&n.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func scanPermissionGroupRow(rows *sql.Rows) (*types.PermissionGroup, error) {
	n := new(types.PermissionGroup)

	err := rows.Scan(
		&n.Id,
		&n.Name,
		&n.Description,
		&n.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func scanResourcePermissionInfoRow(rows *sql.Rows) (*types.GroupResourcePermissionInfo, error) {
	n := new(types.GroupResourcePermissionInfo)

	err := rows.Scan(
		&n.Resource,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func scanActionPermissionInfoRow(rows *sql.Rows) (*types.GroupActionPermissionInfo, error) {
	n := new(types.GroupActionPermissionInfo)

	err := rows.Scan(
		&n.Action,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func buildRoleSearchQuery(
	query types.RolesSearchQuery,
	base string,
) (string, []any) {
	clauses := []string{}
	args := []any{}
	argsPos := 1

	if query.Name != nil {
		clauses = append(clauses, fmt.Sprintf("name ILIKE $%d", argsPos))
		args = append(args, fmt.Sprintf("%%%s%%", *query.Name))
		argsPos++
	}

	q := base
	if len(clauses) > 0 {
		q += " WHERE " + strings.Join(clauses, " AND ")
	}

	q += ";"
	return q, args
}

func buildPermissionGroupSearchQuery(
	query types.PermissionGroupSearchQuery,
	base string,
) (string, []any) {
	clauses := []string{}
	args := []any{}
	argsPos := 1

	if query.Name != nil {
		clauses = append(clauses, fmt.Sprintf("name ILIKE $%d", argsPos))
		args = append(args, fmt.Sprintf("%%%s%%", *query.Name))
		argsPos++
	}

	q := base
	if len(clauses) > 0 {
		q += " WHERE " + strings.Join(clauses, " AND ")
	}

	q += ";"
	return q, args
}
