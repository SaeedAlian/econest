ALTER TABLE addresses
  DROP CONSTRAINT IF EXISTS address_has_owner;

ALTER TABLE addresses
  DROP COLUMN IF EXISTS store_id;

ALTER TABLE addresses
  ALTER COLUMN user_id SET NOT NULL;

ALTER TABLE phonenumbers
  DROP CONSTRAINT IF EXISTS phonenumber_has_owner;

ALTER TABLE phonenumbers
  DROP COLUMN IF EXISTS store_id;

ALTER TABLE phonenumbers
  ALTER COLUMN user_id SET NOT NULL;

DROP TABLE stores_settings;
DROP TABLE store_products;
DROP TABLE stores;
