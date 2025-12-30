-- Rename product 'jaja' to 'jajka'
UPDATE products SET name = 'jajka' WHERE lower(name) = 'jaja';

-- Update shopping list items with the same name
UPDATE shopping_list_items SET name = 'jajka' WHERE lower(name) = 'jaja';
