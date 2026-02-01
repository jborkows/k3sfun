-- Add słoik majonezu and słoik musztardy to przyprawy group

-- Ensure group exists
INSERT OR IGNORE INTO groups(name) VALUES ('przyprawy');

-- Add products (package = opakowanie, start as missing = quantity_value 2 meaning in stock)
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
  ('słoik majonezu', 'majonez', (SELECT id FROM groups WHERE name = 'przyprawy'), 2, 'opakowanie', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('słoik musztardy', 'musztarda', (SELECT id FROM groups WHERE name = 'przyprawy'), 2, 'opakowanie', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
