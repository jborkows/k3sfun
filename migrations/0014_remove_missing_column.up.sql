-- Remove the redundant 'missing' column from products table
-- Missing status is now derived from quantity_value = 0

-- SQLite doesn't support DROP COLUMN directly, so we need to recreate the table

-- Step 1: Rename old table
ALTER TABLE products RENAME TO products_old;

-- Step 2: Create new table without 'missing' column
CREATE TABLE products (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL UNIQUE,
  group_id INTEGER NULL REFERENCES groups(id) ON DELETE SET NULL,
  icon_key TEXT NOT NULL DEFAULT 'cart',
  quantity_value REAL NOT NULL DEFAULT 0,
  quantity_unit TEXT NOT NULL DEFAULT 'sztuk',
  min_quantity_value REAL NOT NULL DEFAULT 0,
  integer_only INTEGER NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Step 3: Copy data (without missing column)
INSERT INTO products (id, name, group_id, icon_key, quantity_value, quantity_unit, min_quantity_value, integer_only, created_at, updated_at)
SELECT id, name, group_id, icon_key, quantity_value, quantity_unit, min_quantity_value, integer_only, created_at, updated_at
FROM products_old;

-- Step 4: Drop old table
DROP TABLE products_old;

-- Step 5: Recreate index
CREATE INDEX IF NOT EXISTS idx_products_group_id ON products(group_id);

-- Step 6: Drop and recreate view without 'missing' column
DROP VIEW IF EXISTS v_products;

CREATE VIEW v_products AS
SELECT
  p.id,
  p.name,
  p.icon_key,
  p.group_id,
  g.name AS group_name,
  p.quantity_value,
  p.quantity_unit,
  p.min_quantity_value,
  p.integer_only,
  p.updated_at
FROM products p
LEFT JOIN groups g ON g.id = p.group_id;
