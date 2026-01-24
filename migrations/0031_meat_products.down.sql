-- Remove meat products added in migration 0031
DELETE FROM products WHERE name IN (
  'kiełbasa (sztuka)',
  'kabanosy (sztuka)',
  'szynka (gramy)'
);
