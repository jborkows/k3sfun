-- Remove product
DELETE FROM products WHERE name IN ('ser parmezan','słoik żurawina','słoik borówka');

-- Remove icon rule
DELETE FROM product_icon_rules WHERE match_substring IN ('parmezan');
