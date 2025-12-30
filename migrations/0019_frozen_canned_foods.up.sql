-- Add new groups for frozen and canned foods
INSERT OR IGNORE INTO groups(name) VALUES ('mrożonki');
INSERT OR IGNORE INTO groups(name) VALUES ('puszki');

-- Add icon rules for frozen foods
INSERT OR IGNORE INTO product_icon_rules(match_substring, icon_key, priority) VALUES
  ('groszek mrożony', 'frozen-peas', 100),
  ('fasolka mrożona', 'frozen-green-beans', 100),
  ('mrożon', 'frozen-peas', 50);

-- Add icon rules for canned goods
INSERT OR IGNORE INTO product_icon_rules(match_substring, icon_key, priority) VALUES
  ('fasola w puszce', 'canned-beans', 100),
  ('kukurydza w puszce', 'canned-corn', 100),
  ('groszek w puszce', 'canned-peas', 100),
  ('w puszce', 'can', 40);

-- Add icon rule for corn on the cob
INSERT OR IGNORE INTO product_icon_rules(match_substring, icon_key, priority) VALUES
  ('kolba kukurydzy', 'corn-cob', 100);

-- Add frozen foods products
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
  ('groszek mrożony', 'frozen-peas', (SELECT id FROM groups WHERE name = 'mrożonki'), 0, 'opakowanie', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('fasolka mrożona', 'frozen-green-beans', (SELECT id FROM groups WHERE name = 'mrożonki'), 0, 'opakowanie', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Add canned goods products
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
  ('czerwona fasola w puszce', 'canned-beans', (SELECT id FROM groups WHERE name = 'puszki'), 0, 'opakowanie', 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('kukurydza w puszce', 'canned-corn', (SELECT id FROM groups WHERE name = 'puszki'), 0, 'opakowanie', 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('fasola biała w puszce', 'canned-beans', (SELECT id FROM groups WHERE name = 'puszki'), 0, 'opakowanie', 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('groszek w puszce', 'canned-peas', (SELECT id FROM groups WHERE name = 'puszki'), 0, 'opakowanie', 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Add corn on the cob to warzywa group
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
  ('kolba kukurydzy', 'corn-cob', (SELECT id FROM groups WHERE name = 'warzywa'), 0, 'sztuk', 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Update papryka to use integer_only
UPDATE products SET integer_only = 1 WHERE name = 'papryka';

-- Update pomidory (plain ones) to use integer_only
UPDATE products SET integer_only = 1 WHERE name = 'pomidory';
