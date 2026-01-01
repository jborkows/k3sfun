-- Revert rename of 'szara reneta'
UPDATE products SET name = 'szara reneta' WHERE lower(name) = 'jabłko szara reneta';

-- Update shopping list items back to old name
UPDATE shopping_list_items SET name = 'szara reneta' WHERE lower(name) = 'jabłko szara reneta';

-- Remove cleaning products
DELETE FROM products WHERE name IN ('ace', 'środek do czyszczenia kuchni', 'środek do lodówek', 'środek do czyszczenia blatów');

-- Remove icon rules
DELETE FROM product_icon_rules WHERE match_substring IN ('ace', 'środek do czyszczenia kuchni', 'środek do lodówek', 'środek do czyszczenia blatów');
