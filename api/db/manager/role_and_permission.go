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
	err := m.db.QueryRow("INSERT INTO roles (name) VALUES ($1) RETURNING id;", p.Name).
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
		q := fmt.Sprintf("SELECT * FROM roles WHERE %s;", strings.Join(clauses, ", "))

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

func (m *Manager) GetRoleWithPermissionGroupsById(id int) (*types.RoleWithPermissionGroups, error) {
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

  rows, err = m.db.Query("SELECT perm", args ...any)

	return role, nil
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

func (m *Manager) GetPermissionGroupWithPermissionsById(
	id int,
) (*types.PermissionGroupWithPermissions, error) {
	// TODO
	// TODO
	// TODO
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

func (m *Manager) AddPermissionGroupToRole(roleId int, permissionGroupId int) error {
	_, err := m.db.Exec(
		"INSERT INTO role_permission_groups (role_id, permission_group_id) VALUES ($1, $2);",
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
		"DELETE FROM role_permission_groups WHERE role_id = $1 AND permission_group_id = $2;",
		roleId,
		permissionGroupId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) AddResourcePermissionToGroup(
	p types.CreateResourcePermissionPayload,
) (int, error) {
	rowId := -1
	err := m.db.QueryRow(
		"INSERT INTO resource_permissions (resource, group_id) VALUES ($1, $2) RETURNING id;",
		p.Resource,
		p.GroupId,
	).Scan(&rowId)
	if err != nil {
		return -1, err
	}

	return rowId, nil
}

func (m *Manager) AddActionPermissionToGroup(p types.CreateActionPermissionPayload) (int, error) {
	rowId := -1
	err := m.db.QueryRow(
		"INSERT INTO action_permissions (action, group_id) VALUES ($1, $2) RETURNING id;",
		p.Action,
		p.GroupId,
	).Scan(&rowId)
	if err != nil {
		return -1, err
	}

	return rowId, nil
}

func (m *Manager) AddResourcePermissionToGroup(
	p types.CreateResourcePermissionPayload,
) (int, error) {
	rowId := -1
	err := m.db.QueryRow(
		"INSERT INTO resource_permissions (resource, group_id) VALUES ($1, $2) RETURNING id;",
		p.Resource,
		p.GroupId,
	).Scan(&rowId)
	if err != nil {
		return -1, err
	}

	return rowId, nil
}

func (m *Manager) AddActionPermissionToGroup(p types.CreateActionPermissionPayload) (int, error) {
	rowId := -1
	err := m.db.QueryRow(
		"INSERT INTO action_permissions (action, group_id) VALUES ($1, $2) RETURNING id;",
		p.Action,
		p.GroupId,
	).Scan(&rowId)
	if err != nil {
		return -1, err
	}

	return rowId, nil
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

func scanRoleRow(rows *sql.Rows) (*types.Role, error) {
	n := new(types.Role)

	err := rows.Scan(
		&n.Id,
		&n.Name,
		&n.CreatedAt,
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

func scanResourcePermission(rows *sql.Rows) (*types.ResourcePermission, error) {
	n := new(types.ResourcePermission)

	err := rows.Scan(
		&n.Id,
		&n.Resource,
		&n.CreatedAt,
		&n.GroupId,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func scanActionPermission(rows *sql.Rows) (*types.ActionPermission, error) {
	n := new(types.ActionPermission)

	err := rows.Scan(
		&n.Id,
		&n.Action,
		&n.CreatedAt,
		&n.GroupId,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}
