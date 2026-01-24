-- Remove lemons product and icon rule
DELETE FROM products WHERE name IN ('cytryny');

DELETE FROM product_icon_rules WHERE match_substring IN ('cytryn');
