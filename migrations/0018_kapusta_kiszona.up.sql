-- Add icon rule for sauerkraut detection
INSERT INTO product_icon_rules(match_substring, icon_key, priority)
VALUES
  ('kapust', 'sauerkraut', 100);

-- Add sauerkraut product to warzywa group
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
  ('kapusta kiszona', 'sauerkraut', (SELECT id FROM groups WHERE name = 'warzywa'), 0, 'kg', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
