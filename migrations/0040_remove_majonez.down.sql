-- Restore majonez product
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
  ('majonez', 'majonez', (SELECT id FROM groups WHERE name = 'przyprawy'), 2, 'opakowanie', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
