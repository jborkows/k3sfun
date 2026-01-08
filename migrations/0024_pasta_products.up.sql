-- Add icon rules for pasta varieties
INSERT INTO product_icon_rules(match_substring, icon_key, priority)
VALUES
  ('spaghetti', 'pasta-spaghetti', 120),
  ('penne', 'pasta-penne', 120),
  ('świderki', 'pasta-spiral', 120),
  ('makaron jajeczny', 'pasta-egg', 120),
  ('plastry lasagni', 'pasta-lasagne', 120),
  ('motylki', 'pasta-bowtie', 120),
  ('dla dzieci', 'pasta-kids', 110),
  ('makaron', 'pasta', 100);

-- Add pasta products to makarony group
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
  ('spaghetti', 'pasta-spaghetti', (SELECT id FROM groups WHERE name = 'makarony'), 0, 'opakowanie', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('penne', 'pasta-penne', (SELECT id FROM groups WHERE name = 'makarony'), 0, 'opakowanie', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('świderki', 'pasta-spiral', (SELECT id FROM groups WHERE name = 'makarony'), 0, 'opakowanie', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('makaron jajeczny', 'pasta-egg', (SELECT id FROM groups WHERE name = 'makarony'), 0, 'opakowanie', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('plastry lasagni', 'pasta-lasagne', (SELECT id FROM groups WHERE name = 'makarony'), 0, 'opakowanie', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('motylki', 'pasta-bowtie', (SELECT id FROM groups WHERE name = 'makarony'), 0, 'opakowanie', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('makaron dla dzieci', 'pasta-kids', (SELECT id FROM groups WHERE name = 'makarony'), 0, 'opakowanie', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
