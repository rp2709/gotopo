DROP TABLE IF EXISTS paths;
CREATE TABLE paths(
  pid SERIAL PRIMARY KEY,
  shape geometry(MultiLineStringZ)
);