-- Add przyprawy group
INSERT OR IGNORE INTO groups(name) VALUES ('przyprawy');

-- Add icon rules for butter, pepper, spices, vegetables, and oils
INSERT INTO product_icon_rules(match_substring, icon_key, priority)
VALUES
  ('masło', 'butter', 100),
  ('maslo', 'butter', 90),
  ('papryka', 'bell-pepper', 100),
  ('papryki', 'bell-pepper', 100),
  ('oregano', 'herb', 100),
  ('sól', 'salt', 100),
  ('sol', 'salt', 90),
  ('pieprz', 'pepper', 100),
  ('bazylia', 'herb', 100),
  ('tymianek', 'herb', 100),
  ('papryka mielona', 'spice', 110),
  ('czosnek', 'garlic', 100),
  ('imbir', 'ginger', 100),
  ('burak', 'beetroot', 100),
  ('olej', 'oil', 100),
  ('oliwa', 'oil', 110),
  ('pomidor', 'tomato', 100),
  ('seler', 'celery', 100),
  ('pietruszka', 'parsley', 100),
  ('sałata', 'lettuce', 100),
  ('dynia', 'pumpkin', 100),
  ('jarmuż', 'kale', 100),
  ('szczypiorek', 'chive', 100),
  ('cebula czerwona', 'onion-red', 110),
  ('cebula', 'onion', 100),
  ('szalotka', 'onion', 90),
  ('koperek', 'herb', 100),
  ('kolendra', 'herb', 100),
  ('natka', 'herb', 100),
  ('pieczark', 'mushroom', 100),
  ('gruszk', 'pear', 100);

-- Add butter to nabiał group
INSERT OR IGNORE INTO products (
  name,
  icon_key,
  group_id,
  quantity_value,
  quantity_unit,
  min_quantity_value,
  missing,
  created_at,
  updated_at
)
VALUES
  ('masło', 'butter', (SELECT id FROM groups WHERE name = 'nabiał'), 0, 'sztuk', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Add vegetables to warzywa group
INSERT OR IGNORE INTO products (
  name,
  icon_key,
  group_id,
  quantity_value,
  quantity_unit,
  min_quantity_value,
  missing,
  created_at,
  updated_at
)
VALUES
  ('papryka', 'bell-pepper', (SELECT id FROM groups WHERE name = 'warzywa'), 0, 'sztuk', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('buraki', 'beetroot', (SELECT id FROM groups WHERE name = 'warzywa'), 0, 'kg', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('pomidory', 'tomato', (SELECT id FROM groups WHERE name = 'warzywa'), 0, 'sztuk', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('pomidory na przetwory', 'tomato', (SELECT id FROM groups WHERE name = 'warzywa'), 0, 'kg', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('pomidory koktajlowe', 'tomato', (SELECT id FROM groups WHERE name = 'warzywa'), 0, 'pęczek', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('seler', 'celery', (SELECT id FROM groups WHERE name = 'warzywa'), 0, 'sztuk', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('pietruszka', 'parsley', (SELECT id FROM groups WHERE name = 'warzywa'), 0, 'sztuk', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('seler naciowy', 'celery', (SELECT id FROM groups WHERE name = 'warzywa'), 0, 'sztuk', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('sałata masłowa', 'lettuce', (SELECT id FROM groups WHERE name = 'warzywa'), 0, 'sztuk', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('sałata lodowa', 'lettuce', (SELECT id FROM groups WHERE name = 'warzywa'), 0, 'sztuk', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('sałata rzymska', 'lettuce', (SELECT id FROM groups WHERE name = 'warzywa'), 0, 'sztuk', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('dynia', 'pumpkin', (SELECT id FROM groups WHERE name = 'warzywa'), 0, 'sztuk', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('jarmuż', 'kale', (SELECT id FROM groups WHERE name = 'warzywa'), 0, 'pęczek', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('szczypiorek', 'chive', (SELECT id FROM groups WHERE name = 'warzywa'), 0, 'pęczek', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('cebula', 'onion', (SELECT id FROM groups WHERE name = 'warzywa'), 0, 'kg', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('cebula czerwona', 'onion-red', (SELECT id FROM groups WHERE name = 'warzywa'), 0, 'kg', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('szalotka', 'onion', (SELECT id FROM groups WHERE name = 'warzywa'), 0, 'sztuk', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Add spices to przyprawy group
INSERT OR IGNORE INTO products (
  name,
  icon_key,
  group_id,
  quantity_value,
  quantity_unit,
  min_quantity_value,
  missing,
  created_at,
  updated_at
)
VALUES
  ('oregano suszone', 'herb', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'gramy', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('sól', 'salt', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'kg', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('pieprz', 'pepper', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'gramy', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('bazylia', 'herb', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'gramy', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('bazylia suszona', 'herb', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'gramy', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('tymianek', 'herb', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'gramy', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('tymianek suszony', 'herb', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'gramy', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('papryka mielona słodka', 'spice', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'gramy', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('papryka mielona ostra', 'spice', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'gramy', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('czosnek', 'garlic', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'główki', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('imbir', 'ginger', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'sztuk', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('koperek', 'herb', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'pęczek', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('kolendra', 'herb', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'pęczek', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('natka', 'herb', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'pęczek', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('sos sojowy', 'spice', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'opakowanie', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('ocet balsamiczny', 'spice', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'litr', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('ocet ryżowy', 'spice', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'litr', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('ocet spirytusowy', 'spice', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'litr', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Add oils and mushrooms to spożywcze group
INSERT OR IGNORE INTO products (
  name,
  icon_key,
  group_id,
  quantity_value,
  quantity_unit,
  min_quantity_value,
  missing,
  created_at,
  updated_at
)
VALUES
  ('olej rzepakowy', 'oil', (SELECT id FROM groups WHERE name = 'spożywcze'), 0, 'litr', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('oliwa', 'oil', (SELECT id FROM groups WHERE name = 'spożywcze'), 0, 'litr', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('olej sezamowy', 'oil', (SELECT id FROM groups WHERE name = 'spożywcze'), 0, 'litr', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('pieczarki białe', 'mushroom', (SELECT id FROM groups WHERE name = 'spożywcze'), 0, 'gramy', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('pieczarki brązowe', 'mushroom', (SELECT id FROM groups WHERE name = 'spożywcze'), 0, 'gramy', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Add gruszki to owoce group
INSERT OR IGNORE INTO products (
  name,
  icon_key,
  group_id,
  quantity_value,
  quantity_unit,
  min_quantity_value,
  missing,
  created_at,
  updated_at
)
VALUES
  ('gruszki', 'pear', (SELECT id FROM groups WHERE name = 'owoce'), 0, 'sztuk', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Update śmietana products to use 'opakowanie' unit
UPDATE products SET quantity_unit = 'opakowanie' WHERE name LIKE 'śmietana%' OR name LIKE 'śmietanka%';
