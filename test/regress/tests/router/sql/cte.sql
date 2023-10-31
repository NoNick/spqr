\c spqr-console
CREATE SHARDING RULE r1 COLUMN i;
CREATE KEY RANGE kridi1 from 0 to 11 route to sh1;
CREATE KEY RANGE kridi2 from 11 to 21 route to sh2;

\c regress

CREATE TABLE tbl(id int, i int);

INSERT INTO tbl (id, i) VALUES(1, 1);

INSERT INTO tbl (id, i) VALUES(12, 12);
INSERT INTO tbl (id, i) VALUES(13, 12);

with cte (i) as (values (12), (13)) select * from tbl inner join cte on tbl.i = cte.i;

DROP TABLE tbl;
\c spqr-console
DROP KEY RANGE ALL;
DROP SHARDING RULE ALL;
