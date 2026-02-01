-- Add majonez and musztarda to przyprawy group

-- Ensure group exists
INSERT OR IGNORE INTO groups(name) VALUES ('przyprawy');

-- Add icon rules for auto-detection
INSERT OR IGNORE INTO product_icon_rules(match_substring, icon_key, priority) VALUES
  ('majonez', 'majonez', 100),
  ('musztarda', 'musztarda', 100);

-- Add products
INSERT OR IGNORE INTO products (
  name,
  icon_key,
  group_id,
  quantity_value,
  quantity_unit,
  created_at,
  updated_at
)
VALUES
  ('majonez', 'majonez', (SELECT id FROM groups WHERE name = 'przyprawy'), 2, 'opakowanie', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('musztarda', 'musztarda', (SELECT id FROM groups WHERE name = 'przyprawy'), 2, 'opakowanie', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
