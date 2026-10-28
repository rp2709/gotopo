package main

import (
	"context"
	"gotopo/gotopo"
	"gotopo/quadtree"
	"log"
	"log/slog"

	"github.com/jackc/pgx/v5"
	pgxgeom "github.com/twpayne/pgx-geom"
)

const ConnectionString = "user=postgres dbname=swisstlm3d sslmode=prefer host=localhost password=postgres"

func main() {
	slog.Info("Connecting to database")
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	err = pgxgeom.Register(ctx, conn)
	if err != nil {
		log.Fatal(err)
	}

	queries := gotopo.New(conn)

	slog.Info("Constructing graph")
	graph := ConstructGraph(queries, ctx, quadtree.Point{X: 2540215.68, Y: 1181115.23}, 0.5)

	slog.Info("Separating and locating edges")
	locatedEdges := DoubleEdges(&graph.Edges, ctx, queries)

	slog.Info("Saving graph")
	err = SaveGraph(ctx, conn, ToVertexParams(graph), locatedEdges, graph.Srid)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Done")
}
