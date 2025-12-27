-- Remove fruit products
DELETE FROM products WHERE name IN (
  'pomarańcza',
  'jabłka',
  'szara reneta',
  'brzoskwinie',
  'śliwki',
  'truskawki',
  'maliny',
  'borówki'
);

-- Restore trash bag icon to cart
UPDATE products SET icon_key = 'cart' WHERE name LIKE '%worki na śmieci%' OR name LIKE '%worki na śmiecie%';

-- Remove icon rules for fruits and trash bags
DELETE FROM product_icon_rules WHERE match_substring IN (
  'pomarańcz',
  'jabłk',
  'jabłko',
  'reneta',
  'brzoskwin',
  'śliwk',
  'truskawk',
  'malin',
  'borówk',
  'worki na śmieci',
  'worek na śmieci',
  'worki na śmiecie'
);
