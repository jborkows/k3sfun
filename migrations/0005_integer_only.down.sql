-- Recreate the view without integer_only
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
  p.updated_at
FROM products p
LEFT JOIN groups g ON g.id = p.group_id;

-- Remove the integer_only column
-- Note: SQLite doesn't support DROP COLUMN directly in older versions
-- This is supported in SQLite 3.35.0+ (modernc/sqlite supports it)
ALTER TABLE products DROP COLUMN integer_only;
