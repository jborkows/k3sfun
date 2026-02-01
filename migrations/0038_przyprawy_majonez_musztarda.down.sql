-- Remove products
DELETE FROM products WHERE name IN ('majonez', 'musztarda');

-- Remove icon rules
DELETE FROM product_icon_rules WHERE match_substring IN ('majonez', 'musztarda');

-- Remove group if empty
DELETE FROM groups WHERE name = 'przyprawy' AND NOT EXISTS (
  SELECT 1 FROM products WHERE group_id = groups.id
);
