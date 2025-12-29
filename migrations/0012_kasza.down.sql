-- Remove kasza products
DELETE FROM products WHERE name IN ('kasza jaglana', 'kasza pęczak', 'kasza gryczana');

-- Remove kasza icon rules
DELETE FROM product_icon_rules WHERE match_substring IN ('jaglana', 'pęczak', 'gryczana', 'kasza');

-- Remove kasza group
DELETE FROM groups WHERE name = 'kasza';
