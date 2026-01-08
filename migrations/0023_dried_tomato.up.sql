-- Add icon rule for dried tomatoes
INSERT OR IGNORE INTO product_icon_rules(match_substring, icon_key, priority)
VALUES
  ('pomidory suszone', 'dried-tomato', 110),
  ('pomidor suszon', 'dried-tomato', 105);

-- Add dried tomato product to warzywa group
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
  ('pomidory suszone', 'dried-tomato', (SELECT id FROM groups WHERE name = 'warzywa'), 0, 'gramy', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
