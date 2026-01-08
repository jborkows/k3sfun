-- Remove dried tomato product
DELETE FROM products WHERE name = 'pomidory suszone';

-- Remove icon rules for dried tomatoes
DELETE FROM product_icon_rules WHERE match_substring IN ('pomidory suszone', 'pomidor suszon');
