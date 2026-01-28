-- Update icon_key for existing domowe products and add product_icon_rules

-- Update specific existing products to new icons
UPDATE products
SET icon_key = 'extension-cord',
    updated_at = CURRENT_TIMESTAMP
WHERE name = 'przedłużacz' COLLATE NOCASE;

UPDATE products
SET icon_key = 'insulation-tape',
    updated_at = CURRENT_TIMESTAMP
WHERE name = 'taśma izolacyjna' COLLATE NOCASE;

-- Insert product_icon_rules to help future auto-assignment (idempotent)
INSERT INTO product_icon_rules(match_substring, icon_key, priority)
SELECT 'przedłużacz', 'extension-cord', 200
WHERE NOT EXISTS (
  SELECT 1 FROM product_icon_rules WHERE match_substring = 'przedłużacz' AND icon_key = 'extension-cord'
);

INSERT INTO product_icon_rules(match_substring, icon_key, priority)
SELECT 'przedluz', 'extension-cord', 180
WHERE NOT EXISTS (
  SELECT 1 FROM product_icon_rules WHERE match_substring = 'przedluz' AND icon_key = 'extension-cord'
);

INSERT INTO product_icon_rules(match_substring, icon_key, priority)
SELECT 'taśma', 'insulation-tape', 200
WHERE NOT EXISTS (
  SELECT 1 FROM product_icon_rules WHERE match_substring = 'taśma' AND icon_key = 'insulation-tape'
);

INSERT INTO product_icon_rules(match_substring, icon_key, priority)
SELECT 'tasma', 'insulation-tape', 180
WHERE NOT EXISTS (
  SELECT 1 FROM product_icon_rules WHERE match_substring = 'tasma' AND icon_key = 'insulation-tape'
);
