-- name: GetRoads :many
SELECT
    gid,
    geom
FROM swisstlm3d_strassen_strasse;
