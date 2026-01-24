-- Add icon rule and product for lemons (cytryny)
INSERT OR IGNORE INTO product_icon_rules(match_substring, icon_key, priority)
VALUES
  ('cytryn', 'orange', 100);

INSERT OR IGNORE INTO products (
  name,
  icon_key,
  group_id,
  quantity_value,
  quantity_unit,
  min_quantity_value,
  created_at,
  updated_at
)
VALUES
  ('cytryny', 'orange', (SELECT id FROM groups WHERE name = 'owoce'), 0, 'sztuk', 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
