-- Revert brown mushrooms to regular mushroom icon
UPDATE products SET icon_key = 'mushroom' WHERE name = 'pieczarki brązowe';

-- Remove icon rule for brown mushrooms
DELETE FROM product_icon_rules WHERE match_substring = 'pieczarki brązowe';
