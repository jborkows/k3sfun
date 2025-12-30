-- Add new groups
INSERT INTO groups(name) VALUES ('pieczywo');
INSERT INTO groups(name) VALUES ('ryba');

-- Add icon rules for chemia (dishwasher items)
INSERT INTO product_icon_rules(match_substring, icon_key, priority)
VALUES
  ('sól do zmywarki', 'dishwasher-salt', 120),
  ('tabletki do zmywarki', 'dishwasher-tablet', 120),
  ('środek do czyszczenia zmywarki', 'dishwasher-cleaner', 120);

-- Add icon rules for pieczywo (bakery)
INSERT INTO product_icon_rules(match_substring, icon_key, priority)
VALUES
  ('chleb', 'bread', 100),
  ('bagietka', 'baguette', 100),
  ('kajzerka', 'roll', 100),
  ('bułka', 'roll', 90);

-- Add icon rules for spożywcze (baking items)
INSERT INTO product_icon_rules(match_substring, icon_key, priority)
VALUES
  ('proszek do pieczenia', 'baking-powder', 100),
  ('pasta waniliowa', 'vanilla-paste', 100);

-- Add icon rules for owoce (fruits)
INSERT INTO product_icon_rules(match_substring, icon_key, priority)
VALUES
  ('banan', 'banana', 100),
  ('pomarańcz', 'orange', 100),
  ('awokado', 'avocado', 100);

-- Add icon rules for mięso (lard)
INSERT INTO product_icon_rules(match_substring, icon_key, priority)
VALUES
  ('słonina', 'lard', 100);

-- Add icon rules for ryba (fish)
INSERT INTO product_icon_rules(match_substring, icon_key, priority)
VALUES
  ('łosoś', 'salmon', 100),
  ('dorsz', 'cod', 100),
  ('makrela wędzona', 'mackerel', 110),
  ('makrela', 'mackerel', 100),
  ('pasta z makreli', 'fish-paste', 110);

-- Add icon rules for nabiał (dairy)
INSERT INTO product_icon_rules(match_substring, icon_key, priority)
VALUES
  ('bieluch', 'bieluch', 100),
  ('twaróg', 'cottage-cheese', 100);

-- Add chemia products
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
  ('sól do zmywarki', 'dishwasher-salt', (SELECT id FROM groups WHERE name = 'chemia'), 0, 'kg', 0, 1, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('tabletki do zmywarki', 'dishwasher-tablet', (SELECT id FROM groups WHERE name = 'chemia'), 0, 'opakowanie', 0, 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('środek do czyszczenia zmywarki', 'dishwasher-cleaner', (SELECT id FROM groups WHERE name = 'chemia'), 0, 'sztuk', 0, 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Add pieczywo products
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
  ('chleb', 'bread', (SELECT id FROM groups WHERE name = 'pieczywo'), 0, 'sztuk', 0, 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('bagietka', 'baguette', (SELECT id FROM groups WHERE name = 'pieczywo'), 0, 'sztuk', 0, 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('kajzerka', 'roll', (SELECT id FROM groups WHERE name = 'pieczywo'), 0, 'sztuk', 0, 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Add spożywcze products
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
  ('proszek do pieczenia', 'baking-powder', (SELECT id FROM groups WHERE name = 'spożywcze'), 0, 'opakowanie', 0, 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('pasta waniliowa', 'vanilla-paste', (SELECT id FROM groups WHERE name = 'spożywcze'), 0, 'opakowanie', 0, 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Add owoce products
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
  ('banany', 'banana', (SELECT id FROM groups WHERE name = 'owoce'), 0, 'sztuk', 0, 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('pomarańcze', 'orange', (SELECT id FROM groups WHERE name = 'owoce'), 0, 'sztuk', 0, 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('awokado', 'avocado', (SELECT id FROM groups WHERE name = 'owoce'), 0, 'sztuk', 0, 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Add mięso product
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
  ('słonina', 'lard', (SELECT id FROM groups WHERE name = 'mięso'), 0, 'kg', 0, 1, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Add ryba products
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
  ('łosoś', 'salmon', (SELECT id FROM groups WHERE name = 'ryba'), 0, 'kg', 0, 1, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('dorsz', 'cod', (SELECT id FROM groups WHERE name = 'ryba'), 0, 'kg', 0, 1, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('makrela wędzona', 'mackerel', (SELECT id FROM groups WHERE name = 'ryba'), 0, 'kg', 0, 1, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('pasta z makreli', 'fish-paste', (SELECT id FROM groups WHERE name = 'ryba'), 0, 'opakowanie', 0, 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Add nabiał products (if not exist)
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
  ('bieluch', 'bieluch', (SELECT id FROM groups WHERE name = 'nabiał'), 0, 'opakowanie', 0, 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('twaróg do rozrabiania', 'cottage-cheese', (SELECT id FROM groups WHERE name = 'nabiał'), 0, 'opakowanie', 0, 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
