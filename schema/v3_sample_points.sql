CREATE EXTENSION IF NOT EXISTS postgis;

DROP TABLE IF EXISTS sample_point;

CREATE TABLE sample_point(
    sid SERIAL PRIMARY KEY,
    eid INTEGER NOT NULL REFERENCES edges (eid),
    location geometry(Point, 2056) NOT NULL
);

CREATE INDEX sample_point_from_edge ON sample_point (eid); -- we will always JOIN with edges