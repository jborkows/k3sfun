ALTER TABLE units ADD COLUMN singular TEXT NOT NULL DEFAULT '';
ALTER TABLE units ADD COLUMN plural TEXT NOT NULL DEFAULT '';

-- Normalize unit naming for pieces
UPDATE products SET quantity_unit = 'sztuk' WHERE lower(quantity_unit) = 'sztuka';
UPDATE shopping_list_items SET quantity_unit = 'sztuk' WHERE lower(quantity_unit) = 'sztuka';
DELETE FROM units WHERE lower(name) = 'sztuka' AND EXISTS (SELECT 1 FROM units WHERE name = 'sztuk');
UPDATE units SET name = 'sztuk' WHERE lower(name) = 'sztuka';
INSERT OR IGNORE INTO units (name, display_order) VALUES ('sztuk', 1);

-- Default forms to unit name
UPDATE units SET singular = name, plural = name WHERE singular = '' AND plural = '';

-- Override specific pluralization rules
UPDATE units SET singular = 'sztuka', plural = 'sztuk' WHERE name = 'sztuk';
UPDATE units SET singular = 'opakowanie', plural = 'opakowania' WHERE name = 'opakowanie';
UPDATE units SET singular = 'pęczek', plural = 'pęczki' WHERE name = 'pęczek';
