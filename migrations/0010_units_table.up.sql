-- Create units table
CREATE TABLE IF NOT EXISTS units (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  display_order INTEGER NOT NULL DEFAULT 0
);

-- Insert existing units from products (deduplicated)
INSERT OR IGNORE INTO units (name, display_order)
SELECT DISTINCT quantity_unit, 0
FROM products
ORDER BY quantity_unit;

-- Set sensible display order for common units
UPDATE units SET display_order = 1 WHERE name = 'sztuk';
UPDATE units SET display_order = 2 WHERE name = 'kg';
UPDATE units SET display_order = 3 WHERE name = 'gramy';
UPDATE units SET display_order = 4 WHERE name = 'litr';
UPDATE units SET display_order = 5 WHERE name = 'opakowanie';
UPDATE units SET display_order = 6 WHERE name = 'pęczek';
UPDATE units SET display_order = 7 WHERE name = 'główki';
