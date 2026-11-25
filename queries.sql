-- name: GetRoads :many
SELECT objectid              AS gid,
       objektart             AS object_type,
       verkehrsbeschraenkung AS trafic_restriction,
       richtungsgetrennt     AS direction_restriction,
       kreisel               AS roundabout,
       shape
FROM tlm_strasse
WHERE objektart NOT IN (0, 1, 2, 3, 6, 13, 14, 21, 22);
-- see swisstlm3d object catalog

-- name: InsertVertex :copyfrom
INSERT INTO vertices (vid, location)
VALUES ($1, $2);

-- name: InsertEdge :copyfrom
INSERT INTO edges (gid, start_vid, end_vid, shape)
VALUES ($1, $2, $3, $4);

-- name: DeleteAllEdges :exec
DELETE
FROM edges;

-- name: DeleteAllVertices :exec
DELETE
FROM vertices;

-- name: GetAllEdges :many
SELECT eid,
       shape
FROM edges;

-- name: InsertSamplePoint :copyfrom
INSERT INTO sample_point (eid, location)
VALUES ($1, $2);

-- name: InsertSamplePointWithViewshed :copyfrom
INSERT INTO sample_point (eid,location,viewshed)
VALUES($1,$2,$3);

-- name: DeleteSamplePointsWithoutViewshedsWithin :exec
DELETE FROM sample_point
WHERE st_within(location,$1) AND viewshed IS NULL;

-- name: SamplePointsExtent :one
SELECT * FROM sample_point_extent;

-- name: GetSomeSamplePoints :many
SELECT *
FROM sample_point
LIMIT $1;

-- name: GetSamplePointsWithin :many
SELECT *
FROM sample_point
WHERE st_within(location,$1);

-- name: InsertWorkAreas :copyfrom
INSERT INTO work_area (num, shape)
VALUES ($1,$2);

-- name: GetNextWorkArea :one
SELECT *
FROM work_area
WHERE NOT done
ORDER BY num
LIMIT 1;

-- name: UpdateWorkArea :exec
UPDATE work_area SET done = $1 WHERE num = $2;

-- name: UpdateSamplePoint :exec
UPDATE sample_point SET viewshed = $1;

-- name: GetGraphWithin :many
SELECT
    v.vid AS vertex_id,
    e.edge_id,
    e.start_vid,
    e.neighbour_id,
    e.length,
    e.road_type,
    e.surface_type,
    e.viewsheds,
    e.shape
FROM vertices v
INNER JOIN edge_agg e ON e.start_vid = v.vid
WHERE st_within(v.location, $1)
ORDER BY v.vid, e.edge_id;

-- name: GetClosestVertexWithin :one
SELECT vid
FROM vertices
WHERE st_dwithin(location, sqlc.arg(point)::geometry, sqlc.arg(max_distance)::float8)
ORDER BY location <-> sqlc.arg(point)::geometry
LIMIT 1;

-- name: AddPath :exec
INSERT INTO paths (shape)
VALUES(st_linemerge($1));