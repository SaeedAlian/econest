CREATE TABLE product_categories (
  id SERIAL PRIMARY KEY,
  name VARCHAR(127) NOT NULL,
  image_name VARCHAR(255) NOT NULL UNIQUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  parent_category_id INTEGER REFERENCES product_categories(id) ON DELETE CASCADE
);

CREATE TABLE products (
  id SERIAL PRIMARY KEY,
  name VARCHAR(150) NOT NULL,
  slug VARCHAR(255) NOT NULL UNIQUE,
  price FLOAT8 NOT NULL CHECK (price >= 0),
  shipment_factor FLOAT8 NOT NULL CHECK (shipment_factor >= 0 AND shipment_factor <= 1),
  description VARCHAR(4095) NOT NULL DEFAULT '',
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  subcategory_id INTEGER NOT NULL REFERENCES product_categories(id) ON DELETE RESTRICT
);

CREATE TABLE product_offers (
  id SERIAL PRIMARY KEY,
  discount FLOAT8 NOT NULL,
  expire_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  product_id INTEGER NOT NULL UNIQUE REFERENCES products(id) ON DELETE CASCADE
);

CREATE TABLE product_images (
  id SERIAL PRIMARY KEY,
  image_name VARCHAR(255) NOT NULL UNIQUE,
  is_main BOOLEAN NOT NULL DEFAULT FALSE,

  product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE
);

CREATE TABLE product_specs (
  id SERIAL PRIMARY KEY,
  label VARCHAR(255) NOT NULL,
  value VARCHAR(255) NOT NULL,

  product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE
);

CREATE TABLE product_tags (
  id SERIAL PRIMARY KEY,
  name VARCHAR(127) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE product_tag_assignments (
  product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  tag_id INTEGER NOT NULL REFERENCES product_tags(id) ON DELETE CASCADE,
  PRIMARY KEY (product_id, tag_id)
);

CREATE TABLE product_attributes (
  id SERIAL PRIMARY KEY,
  label VARCHAR(255) NOT NULL
);

CREATE TABLE product_attribute_options (
  id SERIAL PRIMARY KEY,
  value VARCHAR(255) NOT NULL,

  attribute_id INTEGER NOT NULL REFERENCES product_attributes(id) ON DELETE CASCADE
);

CREATE TABLE product_variants (
  id SERIAL PRIMARY KEY,
  quantity INTEGER NOT NULL DEFAULT 0 CHECK (quantity >= 0),

  product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE
);

CREATE TABLE product_variant_attribute_options (
  variant_id INTEGER NOT NULL REFERENCES product_variants(id) ON DELETE CASCADE,
  attribute_id INTEGER NOT NULL REFERENCES product_attributes(id) ON DELETE RESTRICT,
  option_id INTEGER NOT NULL REFERENCES product_attribute_options(id) ON DELETE RESTRICT,
  PRIMARY KEY (variant_id, attribute_id)
);

CREATE TABLE product_comments (
  id SERIAL PRIMARY KEY,
  scoring INTEGER NOT NULL CHECK (scoring >= 1 AND scoring <= 5),
  comment VARCHAR(1023),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT
);
