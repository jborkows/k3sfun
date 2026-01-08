-- Remove pasta products
DELETE FROM products WHERE name IN ('spaghetti', 'penne', 'świderki', 'makaron jajeczny', 'plastry lasagni', 'motylki', 'makaron dla dzieci');

-- Remove pasta icon rules
DELETE FROM product_icon_rules WHERE match_substring IN ('spaghetti', 'penne', 'świderki', 'makaron jajeczny', 'plastry lasagni', 'motylki', 'dla dzieci', 'makaron');
