-- Remove the koncentrat product and its icon rule
DELETE FROM products WHERE name = 'koncentrat pomidorowy';

DELETE FROM product_icon_rules WHERE match_substring = 'koncentrat pomidor';

-- Revert 'cytryny' product icon_key back to 'orange' (best-effort) and restore the original rule
UPDATE products
SET icon_key = 'orange',
    updated_at = CURRENT_TIMESTAMP
WHERE name = 'cytryny' COLLATE NOCASE;

DELETE FROM product_icon_rules WHERE match_substring = 'cytryn';
INSERT OR IGNORE INTO product_icon_rules(match_substring, icon_key, priority)
VALUES ('cytryn', 'orange', 100);
