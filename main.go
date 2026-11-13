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

type command int

const (
	commandConstructGraph command = iota
	commandSampleEdges
	viewsheds
)

var commandFromName = map[string]command{
	"construct-graph": commandConstructGraph,
	"sample-edges":    commandSampleEdges,
	"viewsheds":       viewsheds,
}

var commandToName = map[command]string{
	commandConstructGraph: "construct-graph",
	commandSampleEdges:    "sample-edges",
	viewsheds:             "viewsheds",
}

var commandDescription = map[command]string{
	commandConstructGraph: "Construct complete, oriented graph from isolated road polylines",
	commandSampleEdges:    "Sample each edge of the graph at a fixed interval",
	viewsheds:             "Perform actions related to viewsheds, see \"gotopo viewsheds --help\".",
}

type commandFunction func(config gotopo.DatabaseConfig, args []string) error

var commandEntryPoint = map[command]commandFunction{
	commandConstructGraph: graph_construction.Run,
	commandSampleEdges:    sample_edges.Run,
	viewsheds:             viewsheds_generation.Run,
}

func main() {
	flag.Usage = usage

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
	cmd, ok := commandFromName[flag.Arg(0)]
	if !ok {
		flag.Usage()
		os.Exit(2)
	}

	commandArgs := flag.Args()[1:]

	err := commandEntryPoint[cmd](dbConfig, commandArgs)
	if err != nil {
		slog.Error("Gotopo failed", "command", commandToName[cmd], "error", err)
		os.Exit(3 + int(cmd))
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: gotopo [options] <command>")
	fmt.Fprintln(os.Stderr, "Available commands:")
	for name, cmd := range commandFromName {
		fmt.Fprintf(os.Stderr, "  - %s : %s\n", name, commandDescription[cmd])
	}
	fmt.Fprintln(os.Stderr, "Only general options are shown. To view options specific to a command, use :\n", "gotopo <command> --help")
	fmt.Fprintln(os.Stderr, "Every subcommand follows the same pattern")
	fmt.Fprintln(os.Stderr, "Options:")
	flag.PrintDefaults()
}
