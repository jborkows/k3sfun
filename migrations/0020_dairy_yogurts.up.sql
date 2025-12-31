-- Add icon rules for yogurt products
INSERT OR IGNORE INTO product_icon_rules(match_substring, icon_key, priority)
VALUES
  ('jogurt grecki', 'jogurt-grecki', 110),
  ('jogurt', 'jogurt', 100),
  ('kefir', 'kefir', 100);

-- Add yogurt and kefir products to nabiał group
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
  ('kefir', 'kefir', (SELECT id FROM groups WHERE name = 'nabiał'), 0, 'sztuk', 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('jogurt', 'jogurt', (SELECT id FROM groups WHERE name = 'nabiał'), 0, 'sztuk', 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('jogurt grecki', 'jogurt-grecki', (SELECT id FROM groups WHERE name = 'nabiał'), 0, 'sztuk', 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
