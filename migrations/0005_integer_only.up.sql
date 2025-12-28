-- Add integer_only column to products table
ALTER TABLE products ADD COLUMN integer_only INTEGER NOT NULL DEFAULT 0;

-- Update existing products that should be integer-only
-- jaja (eggs), pietruszka (parsley), seler (celery), jarmuż (kale)
-- and products in group 'owoce' (fruits)
UPDATE products SET integer_only = 1
WHERE lower(name) IN ('jaja', 'pietruszka', 'seler', 'jarmuż');

UPDATE products SET integer_only = 1
WHERE group_id = (SELECT id FROM groups WHERE name = 'owoce');

-- Drop and recreate the view to include integer_only
DROP VIEW IF EXISTS v_products;

CREATE VIEW v_products AS
SELECT
  p.id,
  p.name,
  p.icon_key,
  p.group_id,
  g.name AS group_name,
  p.quantity_value,
  p.quantity_unit,
  p.min_quantity_value,
  p.missing,
  p.integer_only,
  p.updated_at
FROM products p
LEFT JOIN groups g ON g.id = p.group_id;
