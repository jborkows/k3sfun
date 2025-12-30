-- Step 6: Drop and recreate view without 'missing' column
DROP VIEW v_products;

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
  p.integer_only,
  p.updated_at
FROM products p
LEFT JOIN groups g ON g.id = p.group_id;
