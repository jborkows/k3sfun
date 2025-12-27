-- Add icon rules for fruits
INSERT INTO product_icon_rules(match_substring, icon_key, priority)
VALUES
  ('pomarańcz', 'orange', 100),
  ('jabłk', 'apple', 100),
  ('jabłko', 'apple', 100),
  ('reneta', 'apple', 110),
  ('brzoskwin', 'peach', 100),
  ('śliwk', 'plum', 100),
  ('truskawk', 'strawberry', 100),
  ('malin', 'raspberry', 100),
  ('borówk', 'blueberry', 100),
  ('worki na śmieci', 'trash-bag', 120),
  ('worek na śmieci', 'trash-bag', 110),
  ('worki na śmiecie', 'trash-bag', 115);

-- Add fruit products
INSERT OR IGNORE INTO products (
  name,
  icon_key,
  group_id,
  quantity_value,
  quantity_unit,
  min_quantity_value,
  missing,
  created_at,
  updated_at
)
VALUES
  ('pomarańcza', 'orange', (SELECT id FROM groups WHERE name = 'owoce'), 0, 'sztuk', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('jabłka', 'apple', (SELECT id FROM groups WHERE name = 'owoce'), 0, 'sztuk', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('szara reneta', 'apple', (SELECT id FROM groups WHERE name = 'owoce'), 0, 'sztuk', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('brzoskwinie', 'peach', (SELECT id FROM groups WHERE name = 'owoce'), 0, 'sztuk', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('śliwki', 'plum', (SELECT id FROM groups WHERE name = 'owoce'), 0, 'sztuk', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('truskawki', 'strawberry', (SELECT id FROM groups WHERE name = 'owoce'), 0, 'kg', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('maliny', 'raspberry', (SELECT id FROM groups WHERE name = 'owoce'), 0, 'opakowanie', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('borówki', 'blueberry', (SELECT id FROM groups WHERE name = 'owoce'), 0, 'kg', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Update trash bag products to use the new icon
UPDATE products SET icon_key = 'trash-bag' WHERE name LIKE '%worki na śmieci%' OR name LIKE '%worki na śmiecie%';
