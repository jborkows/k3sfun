-- Revert product name back to 'jaja'
UPDATE products SET name = 'jaja' WHERE lower(name) = 'jajka';

-- Revert shopping list items
UPDATE shopping_list_items SET name = 'jaja' WHERE lower(name) = 'jajka';
