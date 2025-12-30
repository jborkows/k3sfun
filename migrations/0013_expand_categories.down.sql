-- Remove nabiał products
DELETE FROM products WHERE name IN ('bieluch', 'twaróg do rozrabiania');

-- Remove ryba products
DELETE FROM products WHERE name IN ('łosoś', 'dorsz', 'makrela wędzona', 'pasta z makreli');

-- Remove mięso product
DELETE FROM products WHERE name = 'słonina';

-- Remove owoce products
DELETE FROM products WHERE name IN ('banany', 'pomarańcze', 'awokado');

-- Remove spożywcze products
DELETE FROM products WHERE name IN ('proszek do pieczenia', 'pasta waniliowa');

-- Remove pieczywo products
DELETE FROM products WHERE name IN ('chleb', 'bagietka', 'kajzerka');

-- Remove chemia products
DELETE FROM products WHERE name IN ('sól do zmywarki', 'tabletki do zmywarki', 'środek do czyszczenia zmywarki');

-- Remove icon rules for nabiał
DELETE FROM product_icon_rules WHERE match_substring IN ('bieluch', 'twaróg');

-- Remove icon rules for ryba
DELETE FROM product_icon_rules WHERE match_substring IN ('łosoś', 'dorsz', 'makrela wędzona', 'makrela', 'pasta z makreli');

-- Remove icon rules for mięso
DELETE FROM product_icon_rules WHERE match_substring = 'słonina';

-- Remove icon rules for owoce
DELETE FROM product_icon_rules WHERE match_substring IN ('banan', 'pomarańcz', 'awokado');

-- Remove icon rules for spożywcze
DELETE FROM product_icon_rules WHERE match_substring IN ('proszek do pieczenia', 'pasta waniliowa');

-- Remove icon rules for pieczywo
DELETE FROM product_icon_rules WHERE match_substring IN ('chleb', 'bagietka', 'kajzerka', 'bułka');

-- Remove icon rules for chemia
DELETE FROM product_icon_rules WHERE match_substring IN ('sól do zmywarki', 'tabletki do zmywarki', 'środek do czyszczenia zmywarki');

-- Remove groups
DELETE FROM groups WHERE name = 'ryba';
DELETE FROM groups WHERE name = 'pieczywo';
