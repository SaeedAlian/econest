CREATE TABLE permission_groups (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE,
  description TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE role_group_assignments (
  role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
  permission_group_id INTEGER NOT NULL REFERENCES permission_groups(id) ON DELETE CASCADE,
  PRIMARY KEY (role_id, permission_group_id)
);

CREATE TABLE group_resource_permissions (
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  resource VARCHAR(63) NOT NULL,
  group_id INTEGER NOT NULL REFERENCES permission_groups(id) ON DELETE CASCADE,
  PRIMARY KEY (resource, group_id)
);

ALTER TABLE group_resource_permissions
  ALTER COLUMN resource TYPE resources USING resource::resources;

CREATE TABLE group_action_permissions (
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  action VARCHAR(63) NOT NULL,
  group_id INTEGER NOT NULL REFERENCES permission_groups(id) ON DELETE CASCADE,
  PRIMARY KEY (action, group_id)
);

ALTER TABLE group_action_permissions
  ALTER COLUMN action TYPE actions USING action::actions;
