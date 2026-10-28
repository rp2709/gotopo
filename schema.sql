SET CLIENT_ENCODING TO UTF8;
SET STANDARD_CONFORMING_STRINGS TO ON;
BEGIN;
CREATE TABLE "swisstlm3d_strassen_strasse"
(
    gid          serial,
    "uuid"       varchar(38),
    "datum_aend" date,
    "datum_erst" date,
    "erstell_j"  int4,
    "erstell_m"  int4,
    "revision_j" int4,
    "revision_m" int4,
    "grund_aend" varchar(20),
    "herkunft"   varchar(20),
    "herkunft_j" int4,
    "herkunft_m" int4,
    "objektart"  varchar(50),
    "revision_q" varchar(20),
    "kunstbaute" varchar(30),
    "wanderwege" varchar(20),
    "befahrbark" varchar(10),
    "eroeffnung" date,
    "stufe"      varchar(5),
    "richtungsg" varchar(10),
    "str_name_u" varchar(38),
    "name"       varchar(254),
    "belagsart"  varchar(10),
    "kreisel"    varchar(10),
    "verkehrsbe" varchar(50),
    "eigentueme" varchar(50),
    "verkehrsbd" varchar(50),
    "strname"    varchar(254),
    "geom"       geometry(MultiLineStringZM)
);
ALTER TABLE "swisstlm3d_strassen_strasse"
    ADD PRIMARY KEY (gid);
--SELECT AddGeometryColumn('','swisstlm3d_strassen_strasse','geom','0','MULTILINESTRING',4);

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
    geom      geometry(MultiLineString,2056)

    CHECK (start_vid <> end_vid),
    PRIMARY KEY (gid,start_vid,end_vid)
);

CREATE INDEX vertices_location_gix
    ON vertices USING GIST (location);

CREATE INDEX edges_start_vid_idx
    ON edges (start_vid);

CREATE INDEX edges_end_vid_idx
    ON edges (end_vid);