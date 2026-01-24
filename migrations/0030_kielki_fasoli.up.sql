-- Add kiełki fasoli to group 'warzywa'

-- Ensure group exists
INSERT OR IGNORE INTO groups(name) VALUES ('warzywa');

-- Insert product
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
  ('kiełki fasoli (opakowanie)', 'lettuce', (SELECT id FROM groups WHERE name = 'warzywa'), 0, 'opakowanie', 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
