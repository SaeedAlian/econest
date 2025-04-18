CREATE TYPE "actions" AS ENUM (
  'full_control',

  'can_ban_user',
  'can_unban_user',

  'can_add_store',
  'can_update_store',
  'can_delete_store',

  'can_create_order',
  'can_modify_order',

  'can_create_order_shipment',
  'can_modify_order_shipment',

  'can_define_product_offer',
  'can_modify_product_offer',
  'can_cancel_product_offer',

  'can_add_role',
  'can_delete_role',
  'can_modify_role',

  'can_add_user_with_role',

  'can_add_permission_group',
  'can_delete_permission_group',
  'can_modify_permission_group'

  'can_assign_permission_group_to_role',
  'can_remove_permission_group_from_role',

  'can_add_product_category',
  'can_modify_product_category',
  'can_delete_product_category',

  'can_add_product',
  'can_update_product',
  'can_delete_product',

  'can_create_product_comment',
  'can_update_product_comment',
  'can_delete_product_comment'
);
CREATE TYPE "resources" AS ENUM (
  'full_access',

  'roles_and_permissions', -- access to roles, permission groups, actions and resources permissions

  'users_full_access', -- full user information even if the info is not public
  'users_public_only', -- only public info of user

  'wallet_transactions_full_access', -- full access to all wallet transactions for all users

  'stores_full_access', -- full store information even if the info is not public
  'stores_public_only', -- only public info of store

  'orders_full_access' -- access to all orders with their transactions
);
CREATE TYPE "transaction_types" AS ENUM ('deposit', 'withdraw', 'purchase', 'sale');
CREATE TYPE "transaction_statuses" AS ENUM ('pending', 'successful', 'failed');
CREATE TYPE "order_statuses" AS ENUM ('pending_payment', 'payment_paid', 'cancelled');
CREATE TYPE "shipment_types" AS ENUM ('shipping', 'returning');
CREATE TYPE "shipment_statuses" AS ENUM ('on_the_way', 'delivered', 'cancelled');
