CREATE TABLE permission_groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE role_permission_groups (
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_group_id INTEGER NOT NULL REFERENCES permission_groups(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_group_id)
);

CREATE TABLE resource_permissions (
    id SERIAL PRIMARY KEY,
    resource VARCHAR(63) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    group_id INTEGER NOT NULL REFERENCES permission_groups(id) ON DELETE CASCADE
);

ALTER TABLE resource_permissions
ALTER COLUMN resource TYPE resources USING resource::resources;

CREATE TABLE action_permissions (
    id SERIAL PRIMARY KEY,
    action VARCHAR(63) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    group_id INTEGER NOT NULL REFERENCES permission_groups(id) ON DELETE CASCADE
);

ALTER TABLE action_permissions
ALTER COLUMN action TYPE actions USING action::actions;
