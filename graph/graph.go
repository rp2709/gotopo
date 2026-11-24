package graph

import (
	"github.com/rp2709/gotopo/gotopo"
)

// Vertex is a graph vertex together with its outgoing edges.
type Vertex struct {
	ID    int32
	Edges []gotopo.EdgeAgg
}

// Graph is an adjacency list keyed by vertex id.
type Graph map[int32]Vertex

// FromRows groups the flat rows returned by GetGraphWithin into an adjacency list.
func FromRows(rows []gotopo.GetGraphWithinRow) Graph {
	vertices := make(Graph, len(rows))
	for _, row := range rows {
		v := vertices[row.VertexID]
		v.ID = row.VertexID
		v.Edges = append(v.Edges, gotopo.EdgeAgg{
			EdgeID:      row.EdgeID,
			StartVid:    row.StartVid,
			NeighbourID: row.NeighbourID,
			Length:      row.Length,
			RoadType:    row.RoadType,
			SurfaceType: row.SurfaceType,
			Viewsheds:   row.Viewsheds,
			Shape:       row.Shape,
		})
		vertices[row.VertexID] = v
	}
	return vertices
}
