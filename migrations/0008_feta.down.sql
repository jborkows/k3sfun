-- Remove feta product
DELETE FROM products WHERE name = 'feta';

-- Remove icon rule for feta
DELETE FROM product_icon_rules WHERE match_substring = 'feta';
