-- Update icon rules for vinegars and soy sauce with specific icons
-- Remove old generic rules first
DELETE FROM product_icon_rules WHERE match_substring IN ('ocet', 'sos sojowy');

-- Add specific icon rules for vinegars
INSERT INTO product_icon_rules(match_substring, icon_key, priority) VALUES
  ('ocet balsamiczny', 'vinegar-balsamic', 110),
  ('ocet ryżowy', 'vinegar-rice', 110),
  ('ocet', 'vinegar', 80),
  ('sos sojowy', 'soy-sauce', 100);

-- Update existing products to use new icons
UPDATE products
SET icon_key = 'vinegar-balsamic', updated_at = CURRENT_TIMESTAMP
WHERE lower(name) LIKE '%ocet balsamiczny%';

UPDATE products
SET icon_key = 'vinegar-rice', updated_at = CURRENT_TIMESTAMP
WHERE lower(name) LIKE '%ocet ryżowy%';

UPDATE products
SET icon_key = 'vinegar', updated_at = CURRENT_TIMESTAMP
WHERE lower(name) LIKE '%ocet%' 
  AND lower(name) NOT LIKE '%balsamiczny%' 
  AND lower(name) NOT LIKE '%ryżowy%';

UPDATE products
SET icon_key = 'soy-sauce', updated_at = CURRENT_TIMESTAMP
WHERE lower(name) LIKE '%sos sojowy%';

-- Fix śmietana and feta to use 'opakowanie' unit (in case earlier migration didn't catch them)
UPDATE products
SET quantity_unit = 'opakowanie', integer_only = 1, updated_at = CURRENT_TIMESTAMP
WHERE (lower(name) LIKE '%śmietana%' OR lower(name) LIKE '%śmietanka%' OR lower(name) LIKE '%feta%')
  AND quantity_unit != 'opakowanie';
