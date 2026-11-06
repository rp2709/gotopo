package main

import (
	"flag"
	"fmt"
	"log/slog"
	"math"
	"os"

	"github.com/rp2709/gotopo/gotopo"
	"github.com/rp2709/gotopo/graph_construction"
	"github.com/rp2709/gotopo/sample_edges"
	"github.com/rp2709/gotopo/viewsheds_generation"
)

type Command int

const (
	CommandConstructGraph Command = iota
	CommandSampleEdges
	CommandGenerateViewsheds
)

var CommandFromName = map[string]Command{
	"construct-graph":    CommandConstructGraph,
	"sample-edges":       CommandSampleEdges,
	"generate-viewsheds": CommandGenerateViewsheds,
}

var CommandToName = map[Command]string{
	CommandConstructGraph:    "construct-graph",
	CommandSampleEdges:       "sample-edges",
	CommandGenerateViewsheds: "generate-viewsheds",
}

var CommandDescription = map[Command]string{
	CommandConstructGraph:    "Construct complete, oriented graph from isolated road polylines",
	CommandSampleEdges:       "Sample each edge of the graph at a fixed interval",
	CommandGenerateViewsheds: "Generate a hemispherical viewshed of arbitrary precision every sample point using a digital surface model",
}

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: gotopo [options] <command>")
		fmt.Fprintln(os.Stderr, "Available commands:")
		for name, command := range CommandFromName {
			fmt.Fprintf(os.Stderr, "  - %s : %s\n", name, CommandDescription[command])
		}
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}

	// Database configuration, shared by all commands
	var dbConfig gotopo.DatabaseConfig
	flag.StringVar(&dbConfig.Dbhost, "dbhost", "localhost", "Database host")
	flag.IntVar(&dbConfig.Dbport, "dbport", 5432, "Database port")
	flag.StringVar(&dbConfig.Dbuser, "dbuser", "postgres", "Database user")
	flag.StringVar(&dbConfig.Dbname, "dbname", "default", "Database name")
	flag.StringVar(&dbConfig.Dbpwd, "dbpassword", os.Getenv("DBPASSWORD"), "Database password")

	verbosity := flag.Int("verbosity", int(slog.LevelInfo), "Logging level. 0 -> info, -4 -> debug")

	// graph-construction specific flags
	graphConstructionConfig := graph_construction.GraphConstructionConfig{Db: &dbConfig}
	flag.Int64Var(&graphConstructionConfig.ExplorationLimit, "exploration-limit", math.MaxInt64, "Limits the exploration of the graph to a set amount of vertices")
	flag.Float64Var(&graphConstructionConfig.StartX, "startingX", 2540215.68, "x coordinate of starting point for graph exploration")
	flag.Float64Var(&graphConstructionConfig.StartY, "startingY", 1181115.23, "y coordinate of starting point for graph exploration")

	// sample edges flags
	sampleEdgesConfig := sample_edges.SampleEdgesConfig{Db: &dbConfig}
	flag.Float64Var(&sampleEdgesConfig.SampleInterval, "sample-interval", 10.0, "distance between every sample point in meters")

	// viewshed-generation specific flags
	generateViewshedsConfig := viewsheds_generation.GenerateViewshedsConfig{Db: &dbConfig}
	flag.StringVar(&generateViewshedsConfig.VrtFilePath, "vrt", "", "Path to .vrt file indexing the surface model's tiles")
	flag.Float64Var(&generateViewshedsConfig.MaximumWorkAreaSize, "max-work-area-size", 10000, "maximum size of work areas")
	flag.Float64Var(&generateViewshedsConfig.MaximumDistance, "max-distance", 40000, "maximum distance to explore from each sample point")
	flag.BoolVar(&generateViewshedsConfig.CreateWorkAreas, "create-work-areas", false, "Create work areas")

	flag.Parse()
	slog.SetLogLoggerLevel(slog.Level(*verbosity))

	// read command and branch out accordingly
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}
	command, ok := CommandFromName[flag.Arg(0)]
	if !ok {
		flag.Usage()
		os.Exit(2)
	}

	var err error
	switch command {
	case CommandConstructGraph:
		err = graph_construction.Run(graphConstructionConfig)
	case CommandSampleEdges:
		err = sample_edges.Run(sampleEdgesConfig)
	case CommandGenerateViewsheds:
		err = viewsheds_generation.Run(generateViewshedsConfig)
	}
	if err != nil {
		slog.Error(fmt.Sprintf("command %s failed", CommandToName[command]), "error", err)
		os.Exit(3 + int(command))
	}
}
