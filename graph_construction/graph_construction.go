package graph_construction

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"math"

	"github.com/rp2709/gotopo/gotopo"
	"github.com/rp2709/gotopo/quadtree"

	"github.com/jackc/pgx/v5"
	pgxgeom "github.com/twpayne/pgx-geom"
)

func Run(dbConfig gotopo.DatabaseConfig, args []string) error {

	// graph-construction specific flags
	var graphConstructionConfig GraphConstructionConfig
	graphConstructionFlags := flag.NewFlagSet("graph-construction", flag.ExitOnError)
	graphConstructionFlags.Int64Var(&graphConstructionConfig.ExplorationLimit, "exploration-limit", math.MaxInt64, "Limits the exploration of the graph to a set amount of vertices")
	graphConstructionFlags.Float64Var(&graphConstructionConfig.StartX, "startingX", 2540215.68, "x coordinate of starting point for graph exploration")
	graphConstructionFlags.Float64Var(&graphConstructionConfig.StartY, "startingY", 1181115.23, "y coordinate of starting point for graph exploration")

	err := graphConstructionFlags.Parse(args)
	if err != nil {
		return err
	}

	slog.Info("Connecting to database")
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbConfig.ConnectionString())
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	err = pgxgeom.Register(ctx, conn)
	if err != nil {
		return err
	}

	queries := gotopo.New(conn)

	slog.Info("Loading all roads")
	roadsBuffer, grErr := queries.GetRoads(ctx)
	if grErr != nil {
		log.Fatal(grErr)
	}

	slog.Info("Constructing graph")
	graph := ConstructGraph(roadsBuffer, quadtree.Point{X: graphConstructionConfig.StartX, Y: graphConstructionConfig.StartY}, 0.5, graphConstructionConfig.ExplorationLimit)
	slog.Debug("Graph created with", "n", len(graph.Vertices), "m", len(graph.Edges))

	slog.Info("Separating and locating edges")
	graph.Edges = DoubleEdges(graph.Edges, roadsBuffer)
	slog.Debug("After separation: ", "m", len(graph.Edges))

	slog.Info("Saving graph")
	err = SaveGraph(ctx, conn, graph)
	if err != nil {
		return err
	}

	slog.Info("Done")
	return nil
}

type GraphConstructionConfig struct {
	ExplorationLimit int64
	StartX, StartY   float64
}
