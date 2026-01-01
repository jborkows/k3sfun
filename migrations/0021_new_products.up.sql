-- Add new products: mandarynki, ręczniki papierowe, bułka tarta
-- Consolidate orange products: pomarańcza/pomarańcz/pomarańcze -> pomarańcze

-- Icon rules for new products
INSERT OR IGNORE INTO product_icon_rules(match_substring, icon_key, priority)
VALUES
  ('mandarynk', 'mandarin', 100),
  ('ręczniki papierowe', 'paper-towels', 100),
  ('bułka tarta', 'breadcrumbs', 100);

-- Add mandarynki to owoce
INSERT OR IGNORE INTO products (
  name, icon_key, group_id, quantity_value, quantity_unit, min_quantity_value, integer_only, created_at, updated_at
)
VALUES
  ('mandarynki', 'mandarin', (SELECT id FROM groups WHERE name = 'owoce'), 0, 'sztuk', 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Add ręczniki papierowe to domowe
INSERT OR IGNORE INTO products (
  name, icon_key, group_id, quantity_value, quantity_unit, min_quantity_value, integer_only, created_at, updated_at
)
VALUES
  ('ręczniki papierowe', 'paper-towels', (SELECT id FROM groups WHERE name = 'domowe'), 0, 'sztuk', 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Add bułka tarta to domowe
INSERT OR IGNORE INTO products (
  name, icon_key, group_id, quantity_value, quantity_unit, min_quantity_value, integer_only, created_at, updated_at
)
VALUES
  ('bułka tarta', 'breadcrumbs', (SELECT id FROM groups WHERE name = 'domowe'), 0, 'sztuk', 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Consolidate orange products: keep pomarańcze, remove duplicates
-- Step 1: Update shopping list items referencing old names to use pomarańcze
UPDATE shopping_list_items 
SET product_id = (SELECT id FROM products WHERE name = 'pomarańcze')
WHERE lower(name) IN ('pomarańcza', 'pomarańcz')
  AND product_id IS NULL;

-- Step 2: Delete shopping list items linked to old product IDs
DELETE FROM shopping_list_items 
WHERE product_id IN (SELECT id FROM products WHERE name IN ('pomarańcza', 'pomarańcz'));

-- Step 3: Delete duplicate products
DELETE FROM products WHERE name IN ('pomarańcza', 'pomarańcz');

-- Step 4: Ensure pomarańcze has the correct icon
UPDATE products SET icon_key = 'orange' WHERE name = 'pomarańcze';
