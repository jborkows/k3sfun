-- Add icon rule for parmezan (Parmesan cheese)
INSERT INTO product_icon_rules(match_substring, icon_key, priority)
VALUES
  ('parmezan', 'parmezan', 120);

-- Add "ser parmezan" product to nabiał group
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
  ('ser parmezan', 'parmezan', (SELECT id FROM groups WHERE name = 'nabiał'), 0, 'opakowanie', 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('słoik żurawina', 'jar-cranberry', (SELECT id FROM groups WHERE name = 'owoce'), 0, 'opakowanie', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('słoik borówka', 'jar-blueberry', (SELECT id FROM groups WHERE name = 'owoce'), 0, 'opakowanie', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
