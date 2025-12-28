-- Revert vinegar and soy sauce icon changes
DELETE FROM product_icon_rules WHERE match_substring IN ('ocet balsamiczny', 'ocet ryżowy', 'ocet', 'sos sojowy');

-- Restore old generic rules
INSERT INTO product_icon_rules(match_substring, icon_key, priority) VALUES
  ('ocet', 'spice', 80),
  ('sos sojowy', 'spice', 100);

-- Revert product icons to spice
UPDATE products
SET icon_key = 'spice', updated_at = CURRENT_TIMESTAMP
WHERE lower(name) LIKE '%ocet%' OR lower(name) LIKE '%sos sojowy%';
