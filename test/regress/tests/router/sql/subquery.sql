\c spqr-console
CREATE SHARDING RULE r1 COLUMN i;
CREATE KEY RANGE kridi1 from 0 to 11 route to sh1;
CREATE KEY RANGE kridi2 from 11 to 21 route to sh2;

\c regress

CREATE TABLE tbl(id int, i int);

INSERT INTO tbl (id, i) VALUES(1, 1);

INSERT INTO tbl (id, i) VALUES(12, 12);
INSERT INTO tbl (id, i) VALUES(13, 12);

delete from tbl
where id in (select id from tbl where i = 12 for update skip locked limit 1) and i = 12 returning *;

DROP TABLE tbl;
\c spqr-console
DROP KEY RANGE ALL;
DROP SHARDING RULE ALL;
