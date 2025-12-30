-- Remove products
DELETE FROM products WHERE name IN (
  'groszek mrożony',
  'fasolka mrożona',
  'czerwona fasola w puszce',
  'kukurydza w puszce',
  'fasola biała w puszce',
  'groszek w puszce',
  'kolba kukurydzy'
);

-- Remove icon rules
DELETE FROM product_icon_rules WHERE match_substring IN (
  'groszek mrożony',
  'fasolka mrożona',
  'mrożon',
  'fasola w puszce',
  'kukurydza w puszce',
  'groszek w puszce',
  'w puszce',
  'kolba kukurydzy'
);

-- Remove groups
DELETE FROM groups WHERE name IN ('mrożonki', 'puszki');

-- Revert papryka and pomidory integer_only flag
UPDATE products SET integer_only = 0 WHERE name = 'papryka';
UPDATE products SET integer_only = 0 WHERE name = 'pomidory';
