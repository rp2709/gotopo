package graph_construction

import (
	"context"
	"log"
	"log/slog"

	"github.com/rp2709/gotopo/gotopo"
	"github.com/rp2709/gotopo/quadtree"

	"github.com/jackc/pgx/v5"
	pgxgeom "github.com/twpayne/pgx-geom"
)

func Run(config GraphConstructionConfig) error {
	slog.Info("Connecting to database")
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, config.Db.ConnectionString())
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
	graph := ConstructGraph(roadsBuffer, quadtree.Point{X: config.StartX, Y: config.StartY}, 0.5, config.ExplorationLimit)
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
	Db *gotopo.DatabaseConfig

	ExplorationLimit int64
	StartX, StartY   float64
}
