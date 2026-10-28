package graph_construction

import (
	"context"
	"fmt"
	"gotopo/gotopo"
	"gotopo/quadtree"
	"log"
	"log/slog"

	"github.com/jackc/pgx/v5"
	pgxgeom "github.com/twpayne/pgx-geom"
)

const connectionStringFmt = "user=%s dbname=%s sslmode=prefer host=%s password=%s port=%d"

func Run(config GraphConstructionConfig) error {
	previousLogLevel := slog.SetLogLoggerLevel(config.LogLevel)
	defer slog.SetLogLoggerLevel(previousLogLevel)

	slog.Info("Connecting to database")
	ctx := context.Background()
	conn, err := pgx.Connect(
		ctx, fmt.Sprintf(
			connectionStringFmt,
			config.Dbuser,
			config.Dbname,
			config.Dbhost,
			config.Dbpwd,
			config.Dbport,
		),
	)
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
	graph := ConstructGraph(roadsBuffer, quadtree.Point{X: config.StartX, Y: config.StartY}, 0.5, config.ExplorationLimit)
	slog.Debug("Graph created with", "n", len(graph.Vertices), "m", len(graph.Edges))

	slog.Info("Separating and locating edges")
	graph.Edges = DoubleEdges(graph.Edges, roadsBuffer)
	slog.Debug("After separation: ", "m", len(graph.Edges))

	slog.Info("Saving graph")
	err = SaveGraph(ctx, conn, graph)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Done")
	return nil
}

type GraphConstructionConfig struct {
	LogLevel slog.Level
	Dbuser   string
	Dbname   string
	Dbpwd    string
	Dbhost   string
	Dbport   int

	ExplorationLimit int64
	StartX, StartY   float64
}
