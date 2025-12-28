-- Add przyprawy (spices) group if it doesn't exist
INSERT OR IGNORE INTO groups(name) VALUES ('przyprawy');

-- Add spice/herb products
INSERT OR IGNORE INTO products (
  name,
  icon_key,
  group_id,
  quantity_value,
  quantity_unit,
  min_quantity_value,
  missing,
  integer_only,
  created_at,
  updated_at
)
VALUES
  ('oregano suszone', 'herb', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'gramy', 0, 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('sól', 'salt', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'kg', 0, 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('pieprz', 'pepper', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'gramy', 0, 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('bazylia', 'herb', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'gramy', 0, 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('bazylia suszona', 'herb', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'gramy', 0, 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('tymianek', 'herb', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'gramy', 0, 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('tymianek suszony', 'herb', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'gramy', 0, 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('papryka mielona słodka', 'spice', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'gramy', 0, 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('papryka mielona ostra', 'spice', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'gramy', 0, 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('czosnek', 'garlic', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'główki', 0, 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('imbir', 'ginger', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'sztuk', 0, 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('koperek', 'herb', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'pęczek', 0, 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('kolendra', 'herb', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'pęczek', 0, 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('natka', 'herb', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'pęczek', 0, 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('sos sojowy', 'spice', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'opakowanie', 0, 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('ocet balsamiczny', 'spice', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'litr', 0, 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('ocet ryżowy', 'spice', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'litr', 0, 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('ocet spirytusowy', 'spice', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'litr', 0, 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
