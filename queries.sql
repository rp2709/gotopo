-- name: GetRoads :many
SELECT
    gid,
    geom
FROM swisstlm3d_strassen_strasse;

-- name: InsertVertex :exec
INSERT INTO vertices (vid, location)
VALUES ($1, $2);

-- name: InsertEdge :exec
INSERT INTO edges (gid, start_vid, end_vid)
VALUES ($1, $2, $3);

-- name: InsertVertices :copyfrom
INSERT INTO vertices (vid, location)
VALUES ($1, $2);

-- name: InsertEdges :copyfrom
INSERT INTO edges (gid, start_vid, end_vid)
VALUES ($1, $2, $3);
