-- name: GetRoads :many
SELECT
    objectid AS gid,
    objektart AS object_type,
    verkehrsbeschraenkung AS trafic_restriction,
    richtungsgetrennt AS direction_restriction,
    kreisel AS roundabout,
    shape
FROM tlm_strasse
WHERE objektart NOT IN (0,1,2,3,6,13,14,21,22); -- see swisstlm3d object catalog

-- name: InsertVertex :exec
INSERT INTO vertices (vid, location)
VALUES ($1, $2);

-- name: InsertEdge :exec
INSERT INTO edges (gid, start_vid, end_vid, shape)
VALUES ($1, $2, $3, $4);

-- name: DeleteAllEdges :exec
DELETE FROM edges;

-- name: DeleteAllVertices :exec
DELETE FROM vertices;

