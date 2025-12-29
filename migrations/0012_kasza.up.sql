-- Add kasza group
INSERT INTO groups(name) VALUES ('kasza');

-- Add icon rules for kasza varieties
INSERT INTO product_icon_rules(match_substring, icon_key, priority)
VALUES
  ('jaglana', 'groats-millet', 110),
  ('pęczak', 'groats-barley', 110),
  ('gryczana', 'groats-buckwheat', 110),
  ('kasza', 'groats', 100);

-- Add kasza products (jaglana, pęczak, gryczana) with kg unit
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
  ('kasza jaglana', 'groats-millet', (SELECT id FROM groups WHERE name = 'kasza'), 0, 'kg', 0, 1, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('kasza pęczak', 'groats-barley', (SELECT id FROM groups WHERE name = 'kasza'), 0, 'kg', 0, 1, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('kasza gryczana', 'groats-buckwheat', (SELECT id FROM groups WHERE name = 'kasza'), 0, 'kg', 0, 1, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
