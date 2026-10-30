package graph_construction

import (
	"context"
	"log/slog"

	"github.com/rp2709/gotopo/gotopo"
	"github.com/rp2709/gotopo/schema"

	"github.com/jackc/pgx/v5"
)

// SaveGraph stores the vertices and edges of the graph in the vertices and edges tables.
// Everything is done in a single transaction, vertices first so the edge foreign keys are satisfied.
// Edges that are incomplete (a missing end) or that are loops (start == end) are skipped.
// Rows are inserted with the COPY protocol since there can be millions of them.
func SaveGraph(ctx context.Context, conn *pgx.Conn, graph Graph) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // no-op once committed

	slog.Debug("Writing schema")
	if _, err = conn.Exec(ctx, schema.GraphSchema()); err != nil {
		return err
	}

	slog.Debug("Deleting existing graph")
	queries := gotopo.New(tx)
	err = queries.DeleteAllEdges(ctx)
	if err != nil {
		return err
	}
	err = queries.DeleteAllVertices(ctx)
	if err != nil {
		return err
	}

	slog.Debug("Adding vertices")
	if graph.Vertices != nil {
		_, err = queries.InsertVertex(ctx, graph.Vertices)
		if err != nil {
			return err
		}
	}

	slog.Debug("Adding edges")
	if graph.Edges != nil {
		_, err = queries.InsertEdge(ctx, graph.Edges)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}
