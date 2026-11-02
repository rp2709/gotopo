package main

import (
	"flag"
	"fmt"
	"log/slog"
	"math"
	"os"

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
		flag.PrintDefaults()
	}

	dbhost := flag.String("dbhost", "localhost", "Database host")
	dbport := flag.Int("dbport", 5432, "Database port")
	dbuser := flag.String("dbuser", "postgres", "Database user")
	dbname := flag.String("dbname", "default", "Database name")
	dbpassword := flag.String("dbpassword", os.Getenv("DBPASSWORD"), "Database password")
	exploration_limit := flag.Int64("exploration-limit", math.MaxInt64, "Limits the exploration of the graph to a set amount of vertices")
	verbosity := flag.Int("verbosity", int(slog.LevelInfo), "Logging level. 0 -> info, -4 -> debug")
	startX := flag.Float64("startingX", 2540215.68, "x coordinate of starting point for graph exploration")
	startY := flag.Float64("startingY", 1181115.23, "y coordinate of starting point for graph exploration")

	sampleInterval := flag.Float64("sample-interval", 10.0, "distance between every sample point in meters")

	vrtFilePath := flag.String("vrt", "", "Path to .vrt file indexing the surface model's tiles")

	flag.Parse()
	logLevel := slog.Level(*verbosity)

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
		err = graph_construction.Run(graph_construction.GraphConstructionConfig{
			LogLevel:         logLevel,
			StartX:           *startX,
			StartY:           *startY,
			ExplorationLimit: *exploration_limit,
			Dbuser:           *dbuser,
			Dbhost:           *dbhost,
			Dbname:           *dbname,
			Dbpwd:            *dbpassword,
			Dbport:           *dbport,
		})
	case CommandSampleEdges:
		err = sample_edges.Run(sample_edges.SampleEdgesConfig{
			LogLevel:       logLevel,
			Dbuser:         *dbuser,
			Dbhost:         *dbhost,
			Dbname:         *dbname,
			Dbpwd:          *dbpassword,
			Dbport:         *dbport,
			SampleInterval: *sampleInterval,
		})
	case CommandGenerateViewsheds:
		err = viewsheds_generation.Run(viewsheds_generation.GenerateViewshedsConfig{
			VrtFilePath: *vrtFilePath,
		})
	}

	if err != nil {
		slog.Error(fmt.Sprintf("command %s failed", CommandToName[command]), "error", err)
		os.Exit(3 + int(command))
	}
}
