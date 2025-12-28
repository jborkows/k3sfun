-- Add feta cheese to nabiał group
INSERT OR IGNORE INTO products (
  name,
  icon_key,
  group_id,
  quantity_value,
  quantity_unit,
  min_quantity_value,
  missing,
  integer_only,
  created_at,
  updated_at
)
VALUES
  ('feta', 'feta', (SELECT id FROM groups WHERE name = 'nabiał'), 0, 'opakowanie', 0, 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Add icon rule for feta
INSERT OR IGNORE INTO product_icon_rules(match_substring, icon_key, priority) VALUES
  ('feta', 'feta', 100);
