-- name: GetRoads :many
SELECT
    objectid AS gid,
    objektart AS object_type,
    verkehrsbeschraenkung AS trafic_restriction,
    richtungsgetrennt AS direction_restriction,
    kreisel AS roundabout,
    shape
FROM tlm_strasse;

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

