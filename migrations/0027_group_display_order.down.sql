ALTER TABLE groups
DROP COLUMN display_order;

DROP VIEW IF EXISTS v_groups;
CREATE VIEW v_groups AS
SELECT id, name
FROM groups;
