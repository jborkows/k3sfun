ALTER TABLE groups
ADD COLUMN display_order INTEGER NOT NULL DEFAULT 0;

UPDATE groups
SET display_order = CASE lower(name)
  WHEN 'warzywa' THEN 10
  WHEN 'jajka' THEN 20
  WHEN 'owoce' THEN 30
  WHEN 'mięso' THEN 40
  WHEN 'nabiał' THEN 50
  WHEN 'spożywcze' THEN 60
  WHEN 'mąki' THEN 70
  WHEN 'domowe' THEN 80
  WHEN 'chemia' THEN 90
  WHEN 'makarony' THEN 100
  WHEN 'ryż' THEN 110
  ELSE 999
END;

DROP VIEW IF EXISTS v_groups;
CREATE VIEW v_groups AS
SELECT id, name, display_order
FROM groups;
