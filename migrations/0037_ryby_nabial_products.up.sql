-- Rename existing 'ryba' group to 'ryby'
UPDATE groups SET name = 'ryby' WHERE name = 'ryba';

-- Add fish and mascarpone products
INSERT OR IGNORE INTO products (
  name, icon_key, group_id, quantity_value, quantity_unit, min_quantity_value, integer_only, created_at, updated_at
)
VALUES
  ('śledzie', 'mackerel', (SELECT id FROM groups WHERE name = 'ryby'), 0, 'gramy', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('łosoś', 'salmon', (SELECT id FROM groups WHERE name = 'ryby'), 0, 'gramy', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('polędwica z dorsza', 'cod', (SELECT id FROM groups WHERE name = 'ryby'), 0, 'gramy', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('filety z dorsza', 'cod', (SELECT id FROM groups WHERE name = 'ryby'), 0, 'gramy', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('plastry z łososia', 'salmon', (SELECT id FROM groups WHERE name = 'ryby'), 0, 'gramy', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('mascarpone', 'mascarpone', (SELECT id FROM groups WHERE name = 'nabiał'), 0, 'opakowanie', 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
