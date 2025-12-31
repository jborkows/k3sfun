-- Remove yogurt and kefir products
DELETE FROM products WHERE name IN ('kefir', 'jogurt', 'jogurt grecki');

-- Remove icon rules
DELETE FROM product_icon_rules WHERE match_substring IN ('jogurt grecki', 'jogurt', 'kefir');
