CREATE EXTENSION IF NOT EXISTS postgis;

DROP TABLE IF EXISTS sample_point CASCADE;

CREATE TABLE sample_point
(
    sid      SERIAL PRIMARY KEY,
    eid      INTEGER               NOT NULL REFERENCES edges (eid),
    location geometry(Point, 2056) NOT NULL
);

CREATE INDEX sample_point_from_edge ON sample_point (eid); -- we will always JOIN with edges
CREATE INDEX sample_point_gist ON sample_point USING gist (location);
-- we need to query sample points by area to pre-compute viewsheds

-- this view exists only to allow precise sqlc type override
DROP VIEW IF EXISTS sample_point_extent;
CREATE VIEW sample_point_extent AS
SELECT st_setsrid(
               st_extent(location),
               st_srid(2056)
       )::geometry(Polygon, 2056) AS extent
FROM sample_point;