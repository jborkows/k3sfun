-- Fix foreign key references after products table recreation
-- SQLite requires table recreation to modify foreign keys

PRAGMA foreign_keys = OFF;

-- Step 1: Rename current table
ALTER TABLE shopping_list_items RENAME TO shopping_list_items_old;

-- Step 2: Create new table with correct foreign key reference
CREATE TABLE shopping_list_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  product_id INTEGER REFERENCES products(id) ON DELETE SET NULL,
  name TEXT NOT NULL,
  quantity_value REAL NOT NULL DEFAULT 1,
  quantity_unit TEXT NOT NULL DEFAULT 'sztuk',
  done INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  done_at TIMESTAMP NULL
);

-- Step 3: Copy data (without done_at - it will be added by ensureShoppingListDoneAt if missing)
INSERT INTO shopping_list_items (id, product_id, name, quantity_value, quantity_unit, done, created_at)
SELECT id, product_id, name, quantity_value, quantity_unit, done, created_at
FROM shopping_list_items_old;

-- Step 4: Drop old table
DROP TABLE shopping_list_items_old;

-- Step 5: Recreate indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_shopping_list_open_product
  ON shopping_list_items(product_id)
  WHERE done = 0 AND product_id IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_shopping_list_open_name
  ON shopping_list_items(lower(name))
  WHERE done = 0 AND product_id IS NULL;

CREATE INDEX IF NOT EXISTS idx_shopping_list_done_at
  ON shopping_list_items(done, done_at);

PRAGMA foreign_keys = ON;
