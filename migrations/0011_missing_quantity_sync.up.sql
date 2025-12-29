-- Datafix: Sync missing flag with quantity
-- If quantity > 0, product is NOT missing
-- If quantity = 0, product IS missing

-- Clear missing flag for products that have quantity > 0
UPDATE products
SET missing = 0, updated_at = CURRENT_TIMESTAMP
WHERE quantity_value > 0 AND missing = 1;

-- Set missing flag for products that have quantity = 0 and are not already marked missing
UPDATE products
SET missing = 1, updated_at = CURRENT_TIMESTAMP
WHERE quantity_value = 0 AND missing = 0;
