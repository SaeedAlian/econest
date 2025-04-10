DELETE FROM role_permission_groups WHERE
  role_id = (SELECT id FROM roles WHERE name = 'Super Admin') AND
  permission_group_id = (SELECT id FROM permission_groups WHERE name = 'Full Control');

DELETE FROM action_permissions
  WHERE action = 'full_control' AND
  group_id = (SELECT id FROM permission_groups WHERE name = 'Full Control');

DELETE FROM resource_permissions
  WHERE resource = 'full_access' AND
  group_id = (SELECT id FROM permission_groups WHERE name = 'Full Control');

DELETE FROM permission_groups WHERE name = 'Full Control';
DELETE FROM roles WHERE name = 'Super Admin';
