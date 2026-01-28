-- Add 'przedłużacz' and 'taśma izolacyjna' products to group 'domowe'

-- Ensure group exists (safe to run multiple times)
INSERT OR IGNORE INTO groups(name) VALUES ('domowe');

-- Insert 'przedłużacz' if missing
INSERT INTO products (
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
SELECT
  'przedłużacz',
  'extension-cord',
  (SELECT id FROM groups WHERE name = 'domowe'),
  1,
  'sztuk',
  0,
  1,
  CURRENT_TIMESTAMP,
  CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM products WHERE name = 'przedłużacz');

-- Insert 'taśma izolacyjna' if missing
INSERT INTO products (
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
SELECT
  'taśma izolacyjna',
  'insulation-tape',
  (SELECT id FROM groups WHERE name = 'domowe'),
  1,
  'sztuk',
  0,
  1,
  CURRENT_TIMESTAMP,
  CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM products WHERE name = 'taśma izolacyjna');
