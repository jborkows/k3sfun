-- Revert śmietana/śmietanka products to "litr" unit, non-integer, and milk icon
UPDATE products
SET quantity_unit = 'litr', integer_only = 0, icon_key = 'milk', updated_at = CURRENT_TIMESTAMP
WHERE lower(name) LIKE 'śmietana%' OR lower(name) LIKE 'śmietanka%';

-- Remove icon rules for śmietana/śmietanka
DELETE FROM product_icon_rules WHERE match_substring IN ('śmietana', 'śmietanka');
