package main

import (
	"flag"
	"fmt"
	"log/slog"
	"math"
	"os"

	"github.com/rp2709/gotopo/graph_construction"
)

type Command int

const (
	CommandConstructGraph Command = iota
)

var CommandFromName = map[string]Command{
	"construct-graph": CommandConstructGraph,
}

var CommandToName = map[Command]string{
	CommandConstructGraph: "construct-graph",
}

var CommandDescription = map[Command]string{
	CommandConstructGraph: "Construct complete, oriented graph from isolated road polylines",
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
	}

	if err != nil {
		slog.Error(fmt.Sprintf("command %s failed", CommandToName[command]), "error", err)
		os.Exit(3 + int(command))
	}
}
