CREATE EXTENSION IF NOT EXISTS postgis;
DROP TABLE IF EXISTS edges;
DROP TABLE IF EXISTS vertices;

CREATE TABLE vertices
(
    vid      INTEGER PRIMARY KEY,
    location geometry(Point, 2056) NOT NULL
);

CREATE TABLE edges
(
    gid       INTEGER NOT NULL,
    start_vid INTEGER NOT NULL REFERENCES vertices (vid),
    end_vid   INTEGER NOT NULL REFERENCES vertices (vid),
    shape      geometry(LineStringZ, 2056)
        CHECK (start_vid <> end_vid),
    PRIMARY KEY (gid, start_vid, end_vid)
);

CREATE INDEX vertices_location_gix
    ON vertices USING GIST (location);

CREATE INDEX edges_start_vid_idx
    ON edges (start_vid);

CREATE INDEX edges_end_vid_idx
    ON edges (end_vid);
