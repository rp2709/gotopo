package schema

import "embed"

//go:embed v2_graph_tables.sql
var schema embed.FS

func GraphSchema() string {
	bytes, err := schema.ReadFile("v2_graph_tables.sql")
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
