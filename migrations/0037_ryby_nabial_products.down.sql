-- Remove added products (excluding łosoś which existed in migration 0013)
DELETE FROM products WHERE name IN (
  'śledzie',
  'polędwica z dorsza',
  'filety z dorsza',
  'plastry z łososia',
  'mascarpone'
);

-- Rename group back to 'ryba'
UPDATE groups SET name = 'ryba' WHERE name = 'ryby';
