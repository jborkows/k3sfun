-- Remove added products
DELETE FROM products WHERE name IN ('słoik majonezu', 'słoik musztardy');

-- Do not remove group itself (it is shared), but keep cleanup logic here if needed in future
