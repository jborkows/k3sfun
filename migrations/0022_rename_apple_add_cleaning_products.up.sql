-- Rename product 'szara reneta' to 'jabłko szara reneta'
UPDATE products SET name = 'jabłko szara reneta' WHERE lower(name) = 'szara reneta';

-- Update shopping list items with the same name
UPDATE shopping_list_items SET name = 'jabłko szara reneta' WHERE lower(name) = 'szara reneta';

-- Add icon rules for new cleaning products
INSERT OR IGNORE INTO product_icon_rules(match_substring, icon_key, priority)
VALUES
  ('ace', 'bleach', 100),
  ('środek do czyszczenia kuchni', 'kitchen-cleaner', 120),
  ('środek do lodówek', 'fridge-cleaner', 120),
  ('środek do czyszczenia blatów', 'countertop-cleaner', 120);

-- Add ace to chemia
INSERT OR IGNORE INTO products (
  name, icon_key, group_id, quantity_value, quantity_unit, min_quantity_value, integer_only, created_at, updated_at
)
VALUES
  ('ace', 'bleach', (SELECT id FROM groups WHERE name = 'chemia'), 0, 'sztuk', 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Add środek do czyszczenia kuchni to chemia
INSERT OR IGNORE INTO products (
  name, icon_key, group_id, quantity_value, quantity_unit, min_quantity_value, integer_only, created_at, updated_at
)
VALUES
  ('środek do czyszczenia kuchni', 'kitchen-cleaner', (SELECT id FROM groups WHERE name = 'chemia'), 0, 'sztuk', 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Add środek do lodówek to chemia
INSERT OR IGNORE INTO products (
  name, icon_key, group_id, quantity_value, quantity_unit, min_quantity_value, integer_only, created_at, updated_at
)
VALUES
  ('środek do lodówek', 'fridge-cleaner', (SELECT id FROM groups WHERE name = 'chemia'), 0, 'sztuk', 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Add środek do czyszczenia blatów to chemia
INSERT OR IGNORE INTO products (
  name, icon_key, group_id, quantity_value, quantity_unit, min_quantity_value, integer_only, created_at, updated_at
)
VALUES
  ('środek do czyszczenia blatów', 'countertop-cleaner', (SELECT id FROM groups WHERE name = 'chemia'), 0, 'sztuk', 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
