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
	slog.SetLogLoggerLevel(slog.LevelDebug)
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

	slog.Info("Loading all roads")
	roadsBuffer, grErr := queries.GetRoads(ctx)
	if grErr != nil {
		log.Fatal(grErr)
	}

	slog.Info("Constructing graph")
	graph := ConstructGraph(roadsBuffer, quadtree.Point{X: 2540215.68, Y: 1181115.23}, 0.5)
	slog.Debug("Graph created with n=", len(graph.Vertices), " and m=", len(graph.Edges))

	slog.Info("Separating and locating edges")
	graph.Edges = DoubleEdges(graph.Edges, roadsBuffer)
	slog.Debug("Edges after separation : ", len(graph.Edges))

	slog.Info("Saving graph")
	err = SaveGraph(ctx, conn, graph)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Done")
}
