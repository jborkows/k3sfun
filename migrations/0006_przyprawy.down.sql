-- Remove spice/herb products added in this migration
DELETE FROM products WHERE name IN (
  'oregano suszone',
  'sól',
  'pieprz',
  'bazylia',
  'bazylia suszona',
  'tymianek',
  'tymianek suszony',
  'papryka mielona słodka',
  'papryka mielona ostra',
  'czosnek',
  'imbir',
  'koperek',
  'kolendra',
  'natka',
  'sos sojowy',
  'ocet balsamiczny',
  'ocet ryżowy',
  'ocet spirytusowy'
);

-- Remove przyprawy group (only if no products reference it)
DELETE FROM groups WHERE name = 'przyprawy' AND NOT EXISTS (
  SELECT 1 FROM products WHERE group_id = groups.id
);
