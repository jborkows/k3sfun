-- Update śmietana/śmietanka products to use "opakowanie" unit, integer_only, and sour-cream icon
UPDATE products
SET quantity_unit = 'opakowanie', integer_only = 1, icon_key = 'sour-cream', updated_at = CURRENT_TIMESTAMP
WHERE lower(name) LIKE 'śmietana%' OR lower(name) LIKE 'śmietanka%';

-- Add icon rule for śmietana/śmietanka
INSERT OR IGNORE INTO product_icon_rules(match_substring, icon_key, priority) VALUES
  ('śmietana', 'sour-cream', 100),
  ('śmietanka', 'sour-cream', 100);
