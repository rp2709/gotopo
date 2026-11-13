package main

import (
	"flag"
	"fmt"
	"log/slog"
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
		fmt.Fprintln(os.Stderr, "Only general options are shown. To view options specific to a command, use :\n", "gotopo <command> --help")
		fmt.Fprintln(os.Stderr, "Every subcommand follows the same pattern")
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

	commandArgs := flag.Args()[1:]

	var err error
	switch command {
	case CommandConstructGraph:
		err = graph_construction.Run(dbConfig, commandArgs)
	case CommandSampleEdges:
		err = sample_edges.Run(dbConfig, commandArgs)
	case CommandGenerateViewsheds:
		err = viewsheds_generation.Run(dbConfig, commandArgs)
	}
	if err != nil {
		slog.Error(fmt.Sprintf("command %s failed", CommandToName[command]), "error", err)
		os.Exit(3 + int(command))
	}
}
