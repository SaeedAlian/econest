INSERT INTO roles (name) VALUES ('Super Admin');
INSERT INTO roles (name) VALUES ('Admin');
INSERT INTO roles (name) VALUES ('Vendor');
INSERT INTO roles (name) VALUES ('Customer');

INSERT INTO permission_groups 
  (name, description) VALUES
  ('Full Control', 'Full access to all resources and actions'),
  ('User Management', 'Can access all users, ban/unban users'),
  ('Store Moderator', 'Can moderate all stores'),
  ('Personal Store Management', 'Can manage personal stores'),
  ('Order Actions', 'Can setup orders and payments'),
  ('Order Moderator', 'Can access & update all orders'),
  ('Withdrawal Transaction Moderator', 'Can approve or cancel withdraw transactions'),
  ('Role & Permission Management', 'Can manage & manipulate roles and permission groups'),
  ('User Registration', 'Can register any user with any roles (except super admin)'),
  ('Product Category Management', 'Can manage & manipulate product categories'),
  ('Product Attribute Management', 'Can manage & manipulate product attributes'),
  ('Product Attribute Management (Restricted)', 'Can only add product attributes'),
  ('Product Management', 'Can manage & manipulate products'),
  ('Product Comment Moderator', 'Can delete product comments');

INSERT INTO group_resource_permissions
  (resource, group_id) VALUES
  ('full_access', (SELECT id FROM permission_groups WHERE name = 'Full Control')),
  ('roles_and_permissions', (SELECT id FROM permission_groups WHERE name = 'Role & Permission Management')),
  ('users_full_access', (SELECT id FROM permission_groups WHERE name = 'User Management')),
  ('wallet_transactions_full_access', (SELECT id FROM permission_groups WHERE name = 'Withdrawal Transaction Moderator')),
  ('stores_full_access', (SELECT id FROM permission_groups WHERE name = 'Store Moderator')),
  ('orders_full_access', (SELECT id FROM permission_groups WHERE name = 'Order Moderator'));

INSERT INTO group_action_permissions
  (action, group_id) VALUES
  ('can_ban_user', (SELECT id FROM permission_groups WHERE name = 'User Management')),
  ('can_unban_user', (SELECT id FROM permission_groups WHERE name = 'User Management')),

  ('can_add_store', (SELECT id FROM permission_groups WHERE name = 'Personal Store Management')),
  ('can_update_store', (SELECT id FROM permission_groups WHERE name = 'Personal Store Management')),
  ('can_delete_store', (SELECT id FROM permission_groups WHERE name = 'Personal Store Management')),

  ('can_create_order', (SELECT id FROM permission_groups WHERE name = 'Order Actions')),
  ('can_delete_order', (SELECT id FROM permission_groups WHERE name = 'Order Actions')),
  ('can_complete_order_payment', (SELECT id FROM permission_groups WHERE name = 'Order Actions')),
  ('can_cancel_order_payment', (SELECT id FROM permission_groups WHERE name = 'Order Actions')),

  ('can_update_order_shipment', (SELECT id FROM permission_groups WHERE name = 'Order Moderator')),
  ('can_cancel_order_payment', (SELECT id FROM permission_groups WHERE name = 'Order Moderator')),
  ('can_delete_order', (SELECT id FROM permission_groups WHERE name = 'Order Moderator')),

  ('can_approve_withdraw_transaction', (SELECT id FROM permission_groups WHERE name = 'Withdrawal Transaction Moderator')),
  ('can_cancel_withdraw_transaction', (SELECT id FROM permission_groups WHERE name = 'Withdrawal Transaction Moderator')),

  ('can_add_role', (SELECT id FROM permission_groups WHERE name = 'Role & Permission Management')),
  ('can_update_role', (SELECT id FROM permission_groups WHERE name = 'Role & Permission Management')),
  ('can_delete_role', (SELECT id FROM permission_groups WHERE name = 'Role & Permission Management')),
  ('can_add_permission_group', (SELECT id FROM permission_groups WHERE name = 'Role & Permission Management')),
  ('can_delete_permission_group', (SELECT id FROM permission_groups WHERE name = 'Role & Permission Management')),
  ('can_update_permission_group', (SELECT id FROM permission_groups WHERE name = 'Role & Permission Management')),
  ('can_assign_permission_group_to_role', (SELECT id FROM permission_groups WHERE name = 'Role & Permission Management')),
  ('can_remove_permission_group_from_role', (SELECT id FROM permission_groups WHERE name = 'Role & Permission Management')),
  ('can_assign_permission_to_group', (SELECT id FROM permission_groups WHERE name = 'Role & Permission Management')),
  ('can_remove_permission_from_group', (SELECT id FROM permission_groups WHERE name = 'Role & Permission Management')),

  ('can_add_user_with_role', (SELECT id FROM permission_groups WHERE name = 'User Registration')),

  ('can_add_product_category', (SELECT id FROM permission_groups WHERE name = 'Product Category Management')),
  ('can_update_product_category', (SELECT id FROM permission_groups WHERE name = 'Product Category Management')),
  ('can_delete_product_category', (SELECT id FROM permission_groups WHERE name = 'Product Category Management')),

  ('can_add_product_attribute', (SELECT id FROM permission_groups WHERE name = 'Product Attribute Management')),
  ('can_update_product_attribute', (SELECT id FROM permission_groups WHERE name = 'Product Attribute Management')),
  ('can_delete_product_attribute', (SELECT id FROM permission_groups WHERE name = 'Product Attribute Management')),

  ('can_add_product_attribute', (SELECT id FROM permission_groups WHERE name = 'Product Attribute Management (Restricted)')),

  ('can_add_product', (SELECT id FROM permission_groups WHERE name = 'Product Management')),
  ('can_update_product', (SELECT id FROM permission_groups WHERE name = 'Product Management')),
  ('can_delete_product', (SELECT id FROM permission_groups WHERE name = 'Product Management')),
  ('can_add_product_tag', (SELECT id FROM permission_groups WHERE name = 'Product Management')),
  ('can_update_product_tag', (SELECT id FROM permission_groups WHERE name = 'Product Management')),
  ('can_delete_product_tag', (SELECT id FROM permission_groups WHERE name = 'Product Management')),
  ('can_add_product_offer', (SELECT id FROM permission_groups WHERE name = 'Product Management')),
  ('can_update_product_offer', (SELECT id FROM permission_groups WHERE name = 'Product Management')),
  ('can_delete_product_offer', (SELECT id FROM permission_groups WHERE name = 'Product Management')),

  ('can_delete_product_comment', (SELECT id FROM permission_groups WHERE name = 'Product Comment Moderator'));


INSERT INTO role_group_assignments 
  (role_id, permission_group_id) VALUES
  (
    (SELECT id FROM roles WHERE name = 'Super Admin'),
    (SELECT id FROM permission_groups WHERE name = 'Full Control')
  ),

  (
    (SELECT id FROM roles WHERE name = 'Admin'),
    (SELECT id FROM permission_groups WHERE name = 'User Management')
  ),

  (
    (SELECT id FROM roles WHERE name = 'Admin'),
    (SELECT id FROM permission_groups WHERE name = 'Store Moderator')
  ),

  (
    (SELECT id FROM roles WHERE name = 'Admin'),
    (SELECT id FROM permission_groups WHERE name = 'Personal Store Management')
  ),

  (
    (SELECT id FROM roles WHERE name = 'Admin'),
    (SELECT id FROM permission_groups WHERE name = 'Order Moderator')
  ),

  (
    (SELECT id FROM roles WHERE name = 'Admin'),
    (SELECT id FROM permission_groups WHERE name = 'Withdrawal Transaction Moderator')
  ),

  (
    (SELECT id FROM roles WHERE name = 'Admin'),
    (SELECT id FROM permission_groups WHERE name = 'Role & Permission Management')
  ),

  (
    (SELECT id FROM roles WHERE name = 'Admin'),
    (SELECT id FROM permission_groups WHERE name = 'User Registration')
  ),

  (
    (SELECT id FROM roles WHERE name = 'Admin'),
    (SELECT id FROM permission_groups WHERE name = 'Product Category Management')
  ),

  (
    (SELECT id FROM roles WHERE name = 'Admin'),
    (SELECT id FROM permission_groups WHERE name = 'Product Attribute Management')
  ),

  (
    (SELECT id FROM roles WHERE name = 'Admin'),
    (SELECT id FROM permission_groups WHERE name = 'Product Management')
  ),

  (
    (SELECT id FROM roles WHERE name = 'Admin'),
    (SELECT id FROM permission_groups WHERE name = 'Product Comment Moderator')
  ),

  (
    (SELECT id FROM roles WHERE name = 'Vendor'),
    (SELECT id FROM permission_groups WHERE name = 'Personal Store Management')
  ),

  (
    (SELECT id FROM roles WHERE name = 'Vendor'),
    (SELECT id FROM permission_groups WHERE name = 'Product Attribute Management (Restricted)')
  ),

  (
    (SELECT id FROM roles WHERE name = 'Vendor'),
    (SELECT id FROM permission_groups WHERE name = 'Product Management')
  ),

  (
    (SELECT id FROM roles WHERE name = 'Customer'),
    (SELECT id FROM permission_groups WHERE name = 'Order Actions')
  );

CREATE OR REPLACE FUNCTION protect_default_roles()
RETURNS TRIGGER AS $$
BEGIN
  IF TG_OP = 'DELETE' THEN
    IF OLD.name IN ('Super Admin', 'Admin', 'Vendor', 'Customer') THEN
      RAISE EXCEPTION 'cannot delete protected role: %', OLD.name;
    END IF;

  ELSIF TG_OP = 'UPDATE' THEN
    IF OLD.name IN ('Super Admin', 'Admin', 'Vendor', 'Customer') AND NEW.name IS DISTINCT FROM OLD.name THEN
      RAISE EXCEPTION 'cannot rename protected role: %', OLD.name;
    END IF;
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_protect_fixed_roles
BEFORE DELETE OR UPDATE ON roles
FOR EACH ROW
EXECUTE FUNCTION protect_default_roles();
