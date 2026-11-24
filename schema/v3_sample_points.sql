CREATE EXTENSION IF NOT EXISTS postgis;

DROP TABLE IF EXISTS sample_point CASCADE;

CREATE TABLE sample_point
(
    sid      SERIAL PRIMARY KEY,
    eid      INTEGER               NOT NULL REFERENCES edges (eid),
    location geometry(Point, 2056) NOT NULL,
    viewshed bit(128)
);

CREATE INDEX sample_point_from_edge ON sample_point (eid); -- we will always JOIN with edges
CREATE INDEX sample_point_gist ON sample_point USING gist (location);
-- we need to query sample points by area to pre-compute viewsheds

-- this view exists only to allow precise sqlc type override
DROP VIEW IF EXISTS sample_point_extent;
CREATE VIEW sample_point_extent AS
SELECT st_setsrid(
               st_extent(location),
               2056
       )::geometry(Polygon, 2056) AS extent
FROM sample_point;

-- view to render edges with viewshed array
DROP VIEW IF EXISTS edge_agg;
CREATE VIEW edge_agg AS
SELECT
    e.eid AS edge_id,
    e.start_vid AS start_vid,
    e.end_vid AS neighbour_id,
    st_length(e.shape) AS length,
    s.objektart AS road_type,
    s.belagsart AS surface_type,
    ARRAY_AGG(sp.viewshed) AS viewsheds
FROM edges e
INNER JOIN tlm_strasse s ON e.gid = s.objectid
INNER JOIN sample_point sp ON sp.eid = e.eid
GROUP BY (e.eid,e.start_vid,e.end_vid,e.shape,s.objektart,s.belagsart);

-- view to render vertices with their edges as lists
DROP MATERIALIZED VIEW IF EXISTS vertex_agg;
CREATE MATERIALIZED VIEW vertex_agg AS
SELECT
    v.vid AS vertex_id,
    v.location AS location,
    ARRAY_AGG(e) AS edges
FROM vertices v
INNER JOIN edge_agg e ON e.start_vid = v.vid
GROUP BY v.vid;