-- Revert the removal of min_quantity_value column

-- Drop the view first
DROP VIEW IF EXISTS v_products;

-- Rename existing table
ALTER TABLE products RENAME TO products_old;

-- Create new table with min_quantity_value
CREATE TABLE IF NOT EXISTS products (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL UNIQUE,
  group_id INTEGER NULL REFERENCES groups(id) ON DELETE SET NULL,
  icon_key TEXT NOT NULL DEFAULT 'cart',
  quantity_value REAL NOT NULL DEFAULT 0,
  quantity_unit TEXT NOT NULL DEFAULT 'sztuk',
  min_quantity_value REAL NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Copy data from old table, defaulting min_quantity_value to 0
INSERT INTO products (id, name, group_id, icon_key, quantity_value, quantity_unit, min_quantity_value, created_at, updated_at)
SELECT id, name, group_id, icon_key, quantity_value, quantity_unit, 0, created_at, updated_at
FROM products_old;

-- Drop old table
DROP TABLE products_old;

-- Recreate the view with min_quantity_value
CREATE VIEW IF NOT EXISTS v_products AS
SELECT
  p.id,
  p.name,
  p.icon_key,
  p.group_id,
  g.name AS group_name,
  p.quantity_value,
  p.quantity_unit,
  p.min_quantity_value,
  p.updated_at
FROM products p
LEFT JOIN groups g ON g.id = p.group_id;
