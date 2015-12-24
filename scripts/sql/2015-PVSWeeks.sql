INSERT INTO years (year) VALUES (2015);

INSERT INTO weeks (year_id, pvs_id, week) VALUES ((SELECT id FROM years WHERE year = 2015), (SELECT id FROM pvs WHERE type = 'A'), 1);
INSERT INTO weeks (year_id, pvs_id, week) VALUES ((SELECT id FROM years WHERE year = 2015), (SELECT id FROM pvs WHERE type = 'A'), 2);
INSERT INTO weeks (year_id, pvs_id, week) VALUES ((SELECT id FROM years WHERE year = 2015), (SELECT id FROM pvs WHERE type = 'A'), 3); 
INSERT INTO weeks (year_id, pvs_id, week) VALUES ((SELECT id FROM years WHERE year = 2015), (SELECT id FROM pvs WHERE type = 'B'), 4); 
INSERT INTO weeks (year_id, pvs_id, week) VALUES ((SELECT id FROM years WHERE year = 2015), (SELECT id FROM pvs WHERE type = 'C'), 5);
INSERT INTO weeks (year_id, pvs_id, week) VALUES ((SELECT id FROM years WHERE year = 2015), (SELECT id FROM pvs WHERE type = 'C'), 6);
INSERT INTO weeks (year_id, pvs_id, week) VALUES ((SELECT id FROM years WHERE year = 2015), (SELECT id FROM pvs WHERE type = 'C'), 7); 
INSERT INTO weeks (year_id, pvs_id, week) VALUES ((SELECT id FROM years WHERE year = 2015), (SELECT id FROM pvs WHERE type = 'C'), 8);
INSERT INTO weeks (year_id, pvs_id, week) VALUES ((SELECT id FROM years WHERE year = 2015), (SELECT id FROM pvs WHERE type = 'D'), 9); 
INSERT INTO weeks (year_id, pvs_id, week) VALUES ((SELECT id FROM years WHERE year = 2015), (SELECT id FROM pvs WHERE type = 'C'), 10);
INSERT INTO weeks (year_id, pvs_id, week) VALUES ((SELECT id FROM years WHERE year = 2015), (SELECT id FROM pvs WHERE type = 'C'), 11); 
INSERT INTO weeks (year_id, pvs_id, week) VALUES ((SELECT id FROM years WHERE year = 2015), (SELECT id FROM pvs WHERE type = 'A'), 12);
INSERT INTO weeks (year_id, pvs_id, week) VALUES ((SELECT id FROM years WHERE year = 2015), (SELECT id FROM pvs WHERE type = 'A'), 13);
INSERT INTO weeks (year_id, pvs_id, week) VALUES ((SELECT id FROM years WHERE year = 2015), (SELECT id FROM pvs WHERE type = 'A'), 14);
INSERT INTO weeks (year_id, pvs_id, week) VALUES ((SELECT id FROM years WHERE year = 2015), (SELECT id FROM pvs WHERE type = 'A'), 15);
INSERT INTO weeks (year_id, pvs_id, week) VALUES ((SELECT id FROM years WHERE year = 2015), (SELECT id FROM pvs WHERE type = 'A'), 16);
INSERT INTO weeks (year_id, pvs_id, week) VALUES ((SELECT id FROM years WHERE year = 2015), (SELECT id FROM pvs WHERE type = 'A'), 17);
