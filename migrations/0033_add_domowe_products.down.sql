-- Remove products added in 0033_add_domowe_products

-- Remove the products we added (safe if they don't exist)
DELETE FROM products WHERE name IN ('przedłużacz', 'taśma izolacyjna');

-- If the 'domowe' group is empty (no products), remove it
DELETE FROM groups
WHERE name = 'domowe'
  AND NOT EXISTS (SELECT 1 FROM products p WHERE p.group_id = groups.id);
