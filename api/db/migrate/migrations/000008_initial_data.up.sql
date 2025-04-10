INSERT INTO roles (name) VALUES ('Super Admin');

INSERT INTO permission_groups 
  (name, description) VALUES
  ('Full Control', 'Full access to all resources and actions');

INSERT INTO resource_permissions
  (resource, group_id) VALUES
  ('full_access', (SELECT id FROM permission_groups WHERE name = 'Full Control'));

INSERT INTO action_permissions
  (action, group_id) VALUES
  ('full_control', (SELECT id FROM permission_groups WHERE name = 'Full Control'));

INSERT INTO role_permission_groups 
  (role_id, permission_group_id) VALUES
  (
    (SELECT id FROM roles WHERE name = 'Super Admin'),
    (SELECT id FROM permission_groups WHERE name = 'Full Control')
  );
