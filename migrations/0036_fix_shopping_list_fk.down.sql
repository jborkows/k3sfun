-- Revert foreign key fix (not recommended, but for rollback purposes)

PRAGMA foreign_keys = OFF;

CREATE TABLE shopping_list_items_old (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  product_id INTEGER REFERENCES "products_old"(id) ON DELETE SET NULL,
  name TEXT NOT NULL,
  quantity_value REAL NOT NULL DEFAULT 1,
  quantity_unit TEXT NOT NULL DEFAULT 'sztuk',
  done INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  done_at TIMESTAMP NULL
);

INSERT INTO shopping_list_items_old
SELECT id, product_id, name, quantity_value, quantity_unit, done, created_at, done_at
FROM shopping_list_items;

DROP TABLE shopping_list_items;

ALTER TABLE shopping_list_items_old RENAME TO shopping_list_items;

CREATE UNIQUE INDEX idx_shopping_list_open_product
  ON shopping_list_items(product_id)
  WHERE done = 0 AND product_id IS NOT NULL;

CREATE UNIQUE INDEX idx_shopping_list_open_name
  ON shopping_list_items(lower(name))
  WHERE done = 0 AND product_id IS NULL;

CREATE INDEX idx_shopping_list_done_at
  ON shopping_list_items(done, done_at);

PRAGMA foreign_keys = ON;
