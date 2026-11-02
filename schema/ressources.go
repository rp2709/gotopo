package schema

import "embed"

//go:embed v2_graph_tables.sql
//go:embed v3_sample_points.sql
var schema embed.FS

func GraphSchema() string {
	bytes, err := schema.ReadFile("v2_graph_tables.sql")
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func SamplePointsSchema() string {
	bytes, err := schema.ReadFile("v3_sample_points.sql")
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
