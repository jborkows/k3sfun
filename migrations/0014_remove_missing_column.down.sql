-- Restore the 'missing' column to products table

-- Step 1: Rename current table
ALTER TABLE products RENAME TO products_new;

-- Step 2: Create table with 'missing' column
CREATE TABLE products (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL UNIQUE,
  group_id INTEGER NULL REFERENCES groups(id) ON DELETE SET NULL,
  icon_key TEXT NOT NULL DEFAULT 'cart',
  quantity_value REAL NOT NULL DEFAULT 0,
  quantity_unit TEXT NOT NULL DEFAULT 'sztuk',
  min_quantity_value REAL NOT NULL DEFAULT 0,
  missing INTEGER NOT NULL DEFAULT 0,
  integer_only INTEGER NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Step 3: Copy data, derive missing from quantity_value
INSERT INTO products (id, name, group_id, icon_key, quantity_value, quantity_unit, min_quantity_value, missing, integer_only, created_at, updated_at)
SELECT id, name, group_id, icon_key, quantity_value, quantity_unit, min_quantity_value,
       CASE WHEN quantity_value = 0 THEN 1 ELSE 0 END,
       integer_only, created_at, updated_at
FROM products_new;

-- Step 4: Drop new table
DROP TABLE products_new;

-- Step 5: Recreate index
CREATE INDEX IF NOT EXISTS idx_products_group_id ON products(group_id);

-- Step 6: Drop and recreate view with 'missing' column
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
  p.missing,
  p.integer_only,
  p.updated_at
FROM products p
LEFT JOIN groups g ON g.id = p.group_id;
