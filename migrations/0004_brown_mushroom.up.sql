-- Add icon rule for brown mushrooms
INSERT INTO product_icon_rules(match_substring, icon_key, priority)
VALUES
  ('pieczarki brązowe', 'mushroom-brown', 110);

-- Update brown mushrooms to use brown icon
UPDATE products SET icon_key = 'mushroom-brown' WHERE name = 'pieczarki brązowe';
