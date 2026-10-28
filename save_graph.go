package main

import (
	"context"
	"gotopo/gotopo"

	"github.com/jackc/pgx/v5"
	"github.com/twpayne/go-geom"
)

// SaveGraph stores the vertices and edges of the graph in the vertices and edges tables.
// Everything is done in a single transaction, vertices first so the edge foreign keys are satisfied.
// Edges that are incomplete (a missing end) or that are loops (start == end) are skipped.
// Rows are inserted with the COPY protocol since there can be millions of them.
func SaveGraph(ctx context.Context, conn *pgx.Conn, vertices *[]gotopo.InsertVertexParams, edges *[]gotopo.InsertEdgeParams, srid int) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // no-op once committed

	queries := gotopo.New(tx)
	err = queries.DeleteAllEdges(ctx)
	if err != nil {
		return err
	}
	err = queries.DeleteAllVertices(ctx)
	if err != nil {
		return err
	}

	if vertices != nil {
		_, err = tx.CopyFrom(ctx,
			pgx.Identifier{"vertices"},
			[]string{"vid", "location"},
			pgx.CopyFromSlice(len(*vertices), func(i int) ([]any, error) {
				vertex := (*vertices)[i]
				return []any{vertex.Vid, pointXY(vertex.Location, srid)}, nil
			}),
		)
		if err != nil {
			return err
		}
	}

	if edges != nil {
		validEdges := make([]*gotopo.InsertEdgeParams, 0, len(*edges))
		for i := range *edges {
			edge := &(*edges)[i]
			if edge.StartVid == none || edge.EndVid == none || edge.StartVid == edge.EndVid {
				continue
			}
			validEdges = append(validEdges, edge)
		}

		_, err = tx.CopyFrom(ctx,
			pgx.Identifier{"edges"},
			[]string{"gid", "start_vid", "end_vid", "geom"},
			pgx.CopyFromSlice(len(validEdges), func(i int) ([]any, error) {
				edge := validEdges[i]
				return []any{edge.Gid, edge.StartVid, edge.EndVid, multiLineStringXY(edge.Geom, srid)}, nil
			}),
		)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

// pointXY returns a 2D copy of the point with the given srid (the vertices column is geometry(Point, srid)).
func pointXY(point geom.Point, srid int) *geom.Point {
	return geom.NewPointFlat(geom.XY, []float64{point.X(), point.Y()}).SetSRID(srid)
}

// multiLineStringXY returns a 2D copy of the geometry with the given srid
// (the edges column is geometry(MultiLineString, srid), whereas the source roads have Z and M values).
// Returns nil (SQL NULL) for an empty geometry.
func multiLineStringXY(multiLineString geom.MultiLineString, srid int) *geom.MultiLineString {
	if multiLineString.NumLineStrings() == 0 {
		return nil
	}
	result := geom.NewMultiLineString(geom.XY).SetSRID(srid)
	for i := 0; i < multiLineString.NumLineStrings(); i++ {
		coords := multiLineString.LineString(i).Coords()
		flat := make([]float64, 0, 2*len(coords))
		for _, coord := range coords {
			flat = append(flat, coord.X(), coord.Y())
		}
		if err := result.Push(geom.NewLineStringFlat(geom.XY, flat)); err != nil {
			panic(err) // layouts always match
		}
	}
	return result
}

// ToVertexParams converts the vertices of the graph to the parameters needed to insert them.
func ToVertexParams(graph *Graph) *[]gotopo.InsertVertexParams {
	params := make([]gotopo.InsertVertexParams, len(graph.Vertices))
	for i, vertex := range graph.Vertices {
		params[i] = gotopo.InsertVertexParams{
			Vid:      vertex.Vid,
			Location: *geom.NewPointFlat(geom.XY, []float64{vertex.Location.X, vertex.Location.Y}),
		}
	}
	return &params
}
