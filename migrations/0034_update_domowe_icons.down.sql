-- Revert icon_key changes made in 0034_update_domowe_icons

-- Revert products to default 'cart' icon if they currently have the updated keys
UPDATE products
SET icon_key = 'cart',
    updated_at = CURRENT_TIMESTAMP
WHERE name = 'przedłużacz' COLLATE NOCASE AND icon_key = 'extension-cord';

UPDATE products
SET icon_key = 'cart',
    updated_at = CURRENT_TIMESTAMP
WHERE name = 'taśma izolacyjna' COLLATE NOCASE AND icon_key = 'insulation-tape';

-- Remove the product_icon_rules we added (only exact matches)
DELETE FROM product_icon_rules WHERE match_substring IN ('przedłużacz','przedluz','taśma','tasma') AND icon_key IN ('extension-cord','insulation-tape');
