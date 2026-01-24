-- Add group (if not present), add icon rule for koncentrat, switch cytryny icon rule to a lemon-specific icon and update existing product
INSERT OR IGNORE INTO groups(name) VALUES ('warzywa');

-- Icon rule for tomato concentrate
INSERT OR IGNORE INTO product_icon_rules(match_substring, icon_key, priority)
VALUES ('koncentrat pomidor', 'tomato', 100);

-- Ensure lemon icon rule for cytryny
DELETE FROM product_icon_rules WHERE match_substring = 'cytryn';
INSERT OR IGNORE INTO product_icon_rules(match_substring, icon_key, priority)
VALUES ('cytryn', 'cytryny', 100);

-- Update any existing product named 'cytryny' to use the new icon_key
UPDATE products
SET icon_key = 'cytryny',
    updated_at = CURRENT_TIMESTAMP
WHERE name = 'cytryny' COLLATE NOCASE;

-- Add the new product (koncentrat pomidorowy) in group 'warzywa'
INSERT OR IGNORE INTO products (
  name,
  icon_key,
  group_id,
  quantity_value,
  quantity_unit,
  min_quantity_value,
  integer_only,
  created_at,
  updated_at
)
VALUES
  ('koncentrat pomidorowy', 'tomato', (SELECT id FROM groups WHERE name = 'warzywa'), 0, 'opakowanie', 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
