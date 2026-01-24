-- Add meat products: kiełbasa, kabanosy, szynka to group 'mięso'

-- Ensure group exists
INSERT OR IGNORE INTO groups(name) VALUES ('mięso');

-- Insert products (use parenthesized names to reflect unit)
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
  ('kiełbasa (sztuka)', 'bacon', (SELECT id FROM groups WHERE name = 'mięso'), 0, 'sztuka', 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('kabanosy (sztuka)', 'bacon', (SELECT id FROM groups WHERE name = 'mięso'), 0, 'sztuka', 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('szynka (gramy)', 'pork', (SELECT id FROM groups WHERE name = 'mięso'), 0, 'gramy', 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
