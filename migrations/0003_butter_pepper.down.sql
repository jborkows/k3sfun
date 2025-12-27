-- Revert śmietana products to 'litr' unit
UPDATE products SET quantity_unit = 'litr' WHERE name LIKE 'śmietana%' OR name LIKE 'śmietanka%';

-- Remove products
DELETE FROM products WHERE name IN (
  'masło', 'papryka', 'buraki',
  'pomidory', 'pomidory na przetwory', 'pomidory koktajlowe',
  'seler', 'pietruszka', 'seler naciowy',
  'sałata masłowa', 'sałata lodowa', 'sałata rzymska',
  'dynia', 'jarmuż', 'szczypiorek',
  'cebula', 'cebula czerwona', 'szalotka',
  'oregano suszone', 'sól', 'pieprz', 'bazylia', 'bazylia suszona',
  'tymianek', 'tymianek suszony',
  'papryka mielona słodka', 'papryka mielona ostra', 'czosnek', 'imbir',
  'koperek', 'kolendra', 'natka',
  'sos sojowy', 'ocet balsamiczny', 'ocet ryżowy', 'ocet spirytusowy',
  'olej rzepakowy', 'oliwa', 'olej sezamowy',
  'pieczarki białe', 'pieczarki brązowe',
  'gruszki'
);

-- Remove icon rules
DELETE FROM product_icon_rules WHERE match_substring IN (
  'masło', 'maslo', 'papryka', 'papryki',
  'oregano', 'sól', 'sol', 'pieprz', 'bazylia', 'tymianek', 'papryka mielona',
  'czosnek', 'imbir', 'burak', 'olej', 'oliwa',
  'pomidor', 'seler', 'pietruszka', 'sałata', 'dynia', 'jarmuż',
  'szczypiorek', 'cebula czerwona', 'cebula', 'szalotka', 'koperek', 'kolendra', 'natka',
  'pieczark', 'gruszk'
);

-- Remove przyprawy group
DELETE FROM groups WHERE name = 'przyprawy';
