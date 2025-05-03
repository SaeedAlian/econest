package db_manager

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

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
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if query.Name != nil {
		clauses = append(clauses, fmt.Sprintf("name ILIKE $%d", argsPos))
		args = append(args, fmt.Sprintf("%%%s%%", *query.Name))
		argsPos++
	}

	var (
		rows *sql.Rows
		err  error
	)

	if len(clauses) == 0 {
		rows, err = m.db.Query("SELECT * FROM roles;")
	} else {
		q := fmt.Sprintf("SELECT * FROM roles WHERE %s;", strings.Join(clauses, " AND "))

		rows, err = m.db.Query(
			q,
			args...,
		)
	}

	if err != nil {
		return nil, err
	}

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

	role := new(types.Role)
	role.Id = -1

	for rows.Next() {
		role, err = scanRoleRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if role.Id == -1 {
		return nil, fmt.Errorf("Role not found")
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

	role := new(types.Role)
	role.Id = -1

	for rows.Next() {
		role, err = scanRoleRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if role.Id == -1 {
		return nil, fmt.Errorf("Role not found")
	}

	return role, nil
}

func (m *Manager) GetPermissionGroups(
	query types.PermissionGroupSearchQuery,
) ([]types.PermissionGroup, error) {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if query.Name != nil {
		clauses = append(clauses, fmt.Sprintf("name ILIKE $%d", argsPos))
		args = append(args, fmt.Sprintf("%%%s%%", *query.Name))
		argsPos++
	}

	var (
		rows *sql.Rows
		err  error
	)

	if len(clauses) == 0 {
		rows, err = m.db.Query("SELECT * FROM permission_groups;")
	} else {
		q := fmt.Sprintf("SELECT * FROM permission_groups WHERE %s;", strings.Join(clauses, " AND "))

		rows, err = m.db.Query(
			q,
			args...,
		)
	}

	if err != nil {
		return nil, err
	}

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

	group := new(types.PermissionGroup)
	group.Id = -1

	for rows.Next() {
		group, err = scanPermissionGroupRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if group.Id == -1 {
		return nil, fmt.Errorf("Permission group not found")
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

	group := new(types.PermissionGroup)
	group.Id = -1

	for rows.Next() {
		group, err = scanPermissionGroupRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if group.Id == -1 {
		return nil, fmt.Errorf("Permission group not found")
	}

	return group, nil
}

func (m *Manager) GetRolesWithPermissionGroups(
	query types.RolesSearchQuery,
) ([]types.RoleWithPermissionGroups, error) {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if query.Name != nil {
		clauses = append(clauses, fmt.Sprintf("name ILIKE $%d", argsPos))
		args = append(args, fmt.Sprintf("%%%s%%", *query.Name))
		argsPos++
	}

	var (
		rows *sql.Rows
		err  error
	)

	if len(clauses) == 0 {
		rows, err = m.db.Query("SELECT * FROM roles;")
	} else {
		q := fmt.Sprintf("SELECT * FROM roles WHERE %s;", strings.Join(clauses, " AND "))

		rows, err = m.db.Query(
			q,
			args...,
		)
	}

	if err != nil {
		return nil, err
	}

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

		groups := []types.PermissionGroup{}

		for groupRows.Next() {
			group, err := scanPermissionGroupRow(groupRows)
			if err != nil {
				return nil, err
			}

			groups = append(groups, *group)
		}

		role_with_groups := types.RoleWithPermissionGroups{
			Id:               role.Id,
			Name:             role.Name,
			Description:      role.Description,
			CreatedAt:        role.CreatedAt,
			UpdatedAt:        role.UpdatedAt,
			PermissionGroups: groups,
		}

		roles = append(roles, role_with_groups)
	}

	return roles, nil
}

func (m *Manager) GetPermissionGroupsWithPermissions(
	query types.PermissionGroupSearchQuery,
) ([]types.PermissionGroupWithPermissions, error) {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if query.Name != nil {
		clauses = append(clauses, fmt.Sprintf("name ILIKE $%d", argsPos))
		args = append(args, fmt.Sprintf("%%%s%%", *query.Name))
		argsPos++
	}

	var (
		rows *sql.Rows
		err  error
	)

	if len(clauses) == 0 {
		rows, err = m.db.Query("SELECT * FROM permission_groups;")
	} else {
		q := fmt.Sprintf("SELECT * FROM permission_groups WHERE %s;", strings.Join(clauses, " AND "))

		rows, err = m.db.Query(
			q,
			args...,
		)
	}

	if err != nil {
		return nil, err
	}

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

		resourceRows, err := m.db.Query(`SELECT
      grp.resource FROM group_resource_permissions grp
      JOIN permission_groups pg ON pg.id = grp.group_id
      WHERE grp.group_id = $1;
    `, group.Id)
		if err != nil {
			return nil, err
		}

		actions := []types.GroupActionPermissionInfo{}
		resources := []types.GroupResourcePermissionInfo{}

		for actionRows.Next() {
			action, err := scanActionPermissionInfo(actionRows)
			if err != nil {
				return nil, err
			}

			actions = append(actions, *action)
		}

		for resourceRows.Next() {
			resource, err := scanResourcePermissionInfo(resourceRows)
			if err != nil {
				return nil, err
			}

			resources = append(resources, *resource)
		}

		group_with_permissions := types.PermissionGroupWithPermissions{
			Id:                  group.Id,
			Name:                group.Name,
			Description:         group.Description,
			CreatedAt:           group.CreatedAt,
			ResourcePermissions: resources,
			ActionPermissions:   actions,
		}

		groups = append(groups, group_with_permissions)
	}

	return groups, nil
}

func (m *Manager) GetRolesBasedOnResourcePermission(resource types.Resource) ([]types.Role, error) {
	rows, err := m.db.Query(`SELECT
    r.* FROM group_resource_permissions grp 
    JOIN role_group_assignments rga ON rga.permission_group_id = grp.group_id
    JOIN roles r ON rga.role_id = r.id WHERE grp.resource = $1;
  `, resource)
	if err != nil {
		return nil, err
	}

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

func (m *Manager) GetRolesBasedOnActionPermission(action types.Action) ([]types.Role, error) {
	rows, err := m.db.Query(`SELECT
    r.* FROM group_action_permissions gap 
    JOIN role_group_assignments rga ON rga.permission_group_id = gap.group_id
    JOIN roles r ON rga.role_id = r.id WHERE gap.action = $1;
  `, action)
	if err != nil {
		return nil, err
	}

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
	resource types.Resource,
) ([]types.PermissionGroup, error) {
	rows, err := m.db.Query(`SELECT
    pg.* FROM group_resource_permissions grp
    JOIN permission_groups pg ON pg.id = grp.group_id
    WHERE grp.resource = $1;
  `, resource)
	if err != nil {
		return nil, err
	}

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
	action types.Action,
) ([]types.PermissionGroup, error) {
	rows, err := m.db.Query(`SELECT
    pg.* FROM group_action_permissions gap
    JOIN permission_groups pg ON pg.id = gap.group_id
    WHERE gap.action = $1;
  `, action)
	if err != nil {
		return nil, err
	}

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

func (m *Manager) AddPermissionGroupToRole(roleId int, permissionGroupId int) error {
	_, err := m.db.Exec(
		"INSERT INTO role_group_assignments (role_id, permission_group_id) VALUES ($1, $2);",
		roleId,
		permissionGroupId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) RemovePermissionGroupFromRole(roleId int, permissionGroupId int) error {
	_, err := m.db.Exec(
		"DELETE FROM role_group_assignments WHERE role_id = $1 AND permission_group_id = $2;",
		roleId,
		permissionGroupId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) AddResourcePermissionToGroup(
	p types.CreateGroupResourcePermissionPayload,
) (int, error) {
	rowId := -1
	err := m.db.QueryRow(
		"INSERT INTO group_resource_permissions (resource, group_id) VALUES ($1, $2) RETURNING id;",
		p.Resource,
		p.GroupId,
	).Scan(&rowId)
	if err != nil {
		return -1, err
	}

	return rowId, nil
}

func (m *Manager) AddActionPermissionToGroup(
	p types.CreateGroupActionPermissionPayload,
) (int, error) {
	rowId := -1
	err := m.db.QueryRow(
		"INSERT INTO group_action_permissions (action, group_id) VALUES ($1, $2) RETURNING id;",
		p.Action,
		p.GroupId,
	).Scan(&rowId)
	if err != nil {
		return -1, err
	}

	return rowId, nil
}

func (m *Manager) RemoveResourcePermissionFromGroup(
	resource types.Resource,
	groupId int,
) error {
	_, err := m.db.Exec(
		"DELETE FROM group_resource_permissions WHERE resource = $1 AND group_id = $2;",
		resource,
		groupId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) RemoveActionPermissionFromGroup(
	action types.Action,
	groupId int,
) error {
	_, err := m.db.Exec(
		"DELETE FROM group_action_permissions WHERE action = $1 AND group_id = $2;",
		action,
		groupId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) UpdateRole(id int, p types.UpdateRolePayload) error {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if p.Name != nil {
		clauses = append(clauses, fmt.Sprintf("name = $%d", argsPos))
		args = append(args, *p.Name)
		argsPos++
	}

	if len(clauses) == 0 {
		return fmt.Errorf("No fields received to update")
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
	args := []interface{}{}
	argsPos := 1

	if p.Name != nil {
		clauses = append(clauses, fmt.Sprintf("name = $%d", argsPos))
		args = append(args, *p.Name)
		argsPos++
	}

	if len(clauses) == 0 {
		return fmt.Errorf("No fields received to update")
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

func scanResourcePermission(rows *sql.Rows) (*types.GroupResourcePermission, error) {
	n := new(types.GroupResourcePermission)

	err := rows.Scan(
		&n.CreatedAt,
		&n.Resource,
		&n.GroupId,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func scanActionPermission(rows *sql.Rows) (*types.GroupActionPermission, error) {
	n := new(types.GroupActionPermission)

	err := rows.Scan(
		&n.CreatedAt,
		&n.Action,
		&n.GroupId,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func scanResourcePermissionInfo(rows *sql.Rows) (*types.GroupResourcePermissionInfo, error) {
	n := new(types.GroupResourcePermissionInfo)

	err := rows.Scan(
		&n.Resource,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func scanActionPermissionInfo(rows *sql.Rows) (*types.GroupActionPermissionInfo, error) {
	n := new(types.GroupActionPermissionInfo)

	err := rows.Scan(
		&n.Action,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}
