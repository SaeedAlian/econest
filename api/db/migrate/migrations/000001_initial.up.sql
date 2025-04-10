CREATE TYPE "actions" AS ENUM (
    'can_add_product',
    'can_update_product',
    'can_delete_product',
    'can_add_vendor',
    'can_update_vendor',
    'can_delete_vendor',
    'can_ban_user',
    'can_unban_user',
    'can_add_product_tag',
    'can_delete_product_tag',
    'can_add_product_category',
    'can_delete_product_category',
    'can_delete_product_comment',
    'can_add_role',
    'can_delete_role',
    'can_modify_role',
    'can_add_permission_group',
    'can_delete_permission_group',
    'can_modify_permission_group'
);
CREATE TYPE "resources" AS ENUM ('products', 'orders', 'users', 'permissions');
CREATE TYPE "transaction_types" AS ENUM ('deposit', 'withdraw', 'purchase', 'sale');
CREATE TYPE "transaction_status" AS ENUM ('pending', 'successful', 'failed');
CREATE TYPE "order_status" AS ENUM ('pending_payment', 'pending_delivery', 'delivered', 'cancelled');
