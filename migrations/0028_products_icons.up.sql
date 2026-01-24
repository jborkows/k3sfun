-- Add groups if missing
INSERT OR IGNORE INTO groups(name) VALUES ('owoce'), ('przyprawy');

-- Add product icon rules (only insert when identical rule doesn't already exist)
INSERT INTO product_icon_rules(match_substring, icon_key, priority)
SELECT 'cytryn', 'cytryny', 100
WHERE NOT EXISTS (SELECT 1 FROM product_icon_rules WHERE match_substring='cytryn' AND icon_key='cytryny');

INSERT INTO product_icon_rules(match_substring, icon_key, priority)
SELECT 'mandaryn', 'mandarynki', 100
WHERE NOT EXISTS (SELECT 1 FROM product_icon_rules WHERE match_substring='mandaryn' AND icon_key='mandarynki');

INSERT INTO product_icon_rules(match_substring, icon_key, priority)
SELECT 'grejfrut', 'grejfrut', 100
WHERE NOT EXISTS (SELECT 1 FROM product_icon_rules WHERE match_substring='grejfrut' AND icon_key='grejfrut');

-- Papryka mielona: prefer exact 'ostra' / 'słodka' matches with higher priority
INSERT INTO product_icon_rules(match_substring, icon_key, priority)
SELECT 'ostra', 'papryka-mielona-ostra', 200
WHERE NOT EXISTS (SELECT 1 FROM product_icon_rules WHERE match_substring='ostra' AND icon_key='papryka-mielona-ostra');

INSERT INTO product_icon_rules(match_substring, icon_key, priority)
SELECT 'słodka', 'papryka-mielona-slodka', 200
WHERE NOT EXISTS (SELECT 1 FROM product_icon_rules WHERE match_substring='słodka' AND icon_key='papryka-mielona-slodka');

-- Also add ascii variant for 'slodka'
INSERT INTO product_icon_rules(match_substring, icon_key, priority)
SELECT 'slodka', 'papryka-mielona-slodka', 200
WHERE NOT EXISTS (SELECT 1 FROM product_icon_rules WHERE match_substring='slodka' AND icon_key='papryka-mielona-slodka');

INSERT INTO product_icon_rules(match_substring, icon_key, priority)
SELECT 'sambal', 'sambal', 100
WHERE NOT EXISTS (SELECT 1 FROM product_icon_rules WHERE match_substring='sambal' AND icon_key='sambal');

INSERT INTO product_icon_rules(match_substring, icon_key, priority)
SELECT 'worcestersh', 'sos-worcestershire', 150
WHERE NOT EXISTS (SELECT 1 FROM product_icon_rules WHERE match_substring='worcestersh' AND icon_key='sos-worcestershire');

INSERT INTO product_icon_rules(match_substring, icon_key, priority)
SELECT 'worcester', 'sos-worcestershire', 150
WHERE NOT EXISTS (SELECT 1 FROM product_icon_rules WHERE match_substring='worcester' AND icon_key='sos-worcestershire');

INSERT INTO product_icon_rules(match_substring, icon_key, priority)
SELECT 'ryb', 'sos-rybny', 150
WHERE NOT EXISTS (SELECT 1 FROM product_icon_rules WHERE match_substring='ryb' AND icon_key='sos-rybny');

-- Add new products
INSERT OR IGNORE INTO products (
  name, icon_key, group_id, quantity_value, quantity_unit, min_quantity_value, integer_only, created_at, updated_at
)
VALUES
  ('grejfrut', 'grejfrut', (SELECT id FROM groups WHERE name = 'owoce'), 0, 'sztuk', 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('sambal', 'sambal', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'opakowanie', 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('sos Worcestershire', 'sos-worcestershire', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'opakowanie', 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('sos rybny', 'sos-rybny', (SELECT id FROM groups WHERE name = 'przyprawy'), 0, 'opakowanie', 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Update pomidory suszone: set unit to 'opakowanie' and replace any existing quantity with 1
UPDATE products
SET quantity_value = 1,
    quantity_unit = 'opakowanie',
    updated_at = CURRENT_TIMESTAMP
WHERE name = 'pomidory suszone' COLLATE NOCASE;

-- Update papryka mielona variants: set unit to 'opakowanie' and set icon keys
UPDATE products
SET quantity_unit = 'opakowanie',
    icon_key = 'papryka-mielona-ostra',
    updated_at = CURRENT_TIMESTAMP
WHERE name = 'papryka mielona ostra' COLLATE NOCASE;

UPDATE products
SET quantity_unit = 'opakowanie',
    icon_key = 'papryka-mielona-slodka',
    updated_at = CURRENT_TIMESTAMP
WHERE name = 'papryka mielona słodka' COLLATE NOCASE OR name = 'papryka mielona slodka' COLLATE NOCASE;

-- Ensure mandarynki receives its new icon by updating existing product if present
UPDATE products
SET icon_key = 'mandarynki',
    updated_at = CURRENT_TIMESTAMP
WHERE name = 'mandarynki' COLLATE NOCASE;
