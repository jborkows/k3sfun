-- Revert new products and orange consolidation

-- Remove new products
DELETE FROM products WHERE name IN ('mandarynki', 'ręczniki papierowe', 'bułka tarta');

-- Remove icon rules
DELETE FROM product_icon_rules WHERE match_substring IN ('mandarynk', 'ręczniki papierowe', 'bułka tarta');

-- Restore pomarańcza (consolidated into pomarańcze)
INSERT OR IGNORE INTO products (
  name, icon_key, group_id, quantity_value, quantity_unit, min_quantity_value, integer_only, created_at, updated_at
)
VALUES
  ('pomarańcza', 'orange', (SELECT id FROM groups WHERE name = 'owoce'), 0, 'sztuk', 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
