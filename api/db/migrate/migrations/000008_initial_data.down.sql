DROP TRIGGER IF EXISTS trg_protect_fixed_roles ON roles;
DROP FUNCTION IF EXISTS protect_default_roles;

DELETE FROM permission_groups WHERE name = 'Full Control';
DELETE FROM permission_groups WHERE name = 'User Management';
DELETE FROM permission_groups WHERE name = 'Store Moderator';
DELETE FROM permission_groups WHERE name = 'Personal Store Management';
DELETE FROM permission_groups WHERE name = 'Order Actions';
DELETE FROM permission_groups WHERE name = 'Order Moderator';
DELETE FROM permission_groups WHERE name = 'Withdrawal Transaction Moderator';
DELETE FROM permission_groups WHERE name = 'Role & Permission Management';
DELETE FROM permission_groups WHERE name = 'User Registration';
DELETE FROM permission_groups WHERE name = 'Product Category Management';
DELETE FROM permission_groups WHERE name = 'Product Attribute Management';
DELETE FROM permission_groups WHERE name = 'Product Attribute Management (Restricted)';
DELETE FROM permission_groups WHERE name = 'Product Management';
DELETE FROM permission_groups WHERE name = 'Product Comment Moderator';

DELETE FROM roles WHERE name = 'Customer';
DELETE FROM roles WHERE name = 'Vendor';
DELETE FROM roles WHERE name = 'Admin';
DELETE FROM roles WHERE name = 'Super Admin';
