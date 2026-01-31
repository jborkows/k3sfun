-- Fix foreign key reference in shopping_list_items table
-- The table was referencing products_old which doesn't exist

-- Disable foreign key checks temporarily
PRAGMA foreign_keys = OFF;

-- Create new table with correct foreign key reference
CREATE TABLE shopping_list_items_new (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  product_id INTEGER REFERENCES products(id) ON DELETE SET NULL,
  name TEXT NOT NULL,
  quantity_value REAL NOT NULL DEFAULT 1,
  quantity_unit TEXT NOT NULL DEFAULT 'sztuk',
  done INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  done_at TIMESTAMP NULL
);

-- Copy data from old table
INSERT INTO shopping_list_items_new
SELECT id, product_id, name, quantity_value, quantity_unit, done, created_at, done_at
FROM shopping_list_items;

-- Drop old table
DROP TABLE shopping_list_items;

-- Rename new table to original name
ALTER TABLE shopping_list_items_new RENAME TO shopping_list_items;

-- Recreate indexes
CREATE UNIQUE INDEX idx_shopping_list_open_product
  ON shopping_list_items(product_id)
  WHERE done = 0 AND product_id IS NOT NULL;

CREATE UNIQUE INDEX idx_shopping_list_open_name
  ON shopping_list_items(lower(name))
  WHERE done = 0 AND product_id IS NULL;

CREATE INDEX idx_shopping_list_done_at
  ON shopping_list_items(done, done_at);

-- Re-enable foreign key checks
PRAGMA foreign_keys = ON;
