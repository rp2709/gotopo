package main

import (
	"context"
	"fmt"
	"gotopo/gotopo"

	"github.com/jackc/pgx/v5"
	"github.com/twpayne/go-geom"
)

// SaveGraph stores the vertices and edges of the graph in the vertices and edges tables.
// Everything is done in a single transaction, vertices first so the edge foreign keys are satisfied.
// Edges that are incomplete (a missing end) or that are loops (start == end) are skipped.
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

	// The generated InsertVertices stores geom.Point by value, but pgx-geom can only encode
	// geom.T, which is implemented by *geom.Point. So copy the vertices with pointers directly.
	rows := make([][]any, len(graph.Vertices))
	for i, v := range graph.Vertices {
		rows[i] = []any{
			v.Vid,
			geom.NewPointFlat(geom.XY, []float64{v.Location.X, v.Location.Y}).SetSRID(srid),
		}
	}
	nv, err := tx.CopyFrom(ctx, pgx.Identifier{"vertices"}, []string{"vid", "location"}, pgx.CopyFromRows(rows))
	if err != nil {
		return err
	}

	ne, err := queries.InsertEdges(ctx, edges)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	fmt.Printf("saved %d vertices and %d edges (%d edges skipped)\n", nv, ne, len(graph.Edges)-len(edges))
	return nil
}
