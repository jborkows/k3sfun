-- Remove sauerkraut product
DELETE FROM products WHERE name = 'kapusta kiszona';

-- Remove icon rule
DELETE FROM product_icon_rules WHERE match_substring = 'kapust';
